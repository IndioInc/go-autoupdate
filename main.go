package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/IndioInc/go-autoupdate/autoupdate"
	"github.com/aws/aws-sdk-go/aws"
)

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("\tgo-autoupdate help")
	fmt.Println("\tgo-autoupdate [flags] init <appName> <channel> <s3Bucket>")
	fmt.Println("\tgo-autoupdate [flags] release <appName> <channel> <s3Bucket> <releasesDir> <releasesTag>")
	fmt.Println("\tgo-autoupdate [flags] download-latest <appName> <channel> <s3Bucket> <outputName>")
	fmt.Println("")
	fmt.Println("Flags:")
	flag.PrintDefaults()
}

func initBucket(awsConfig *aws.Config) {
	appName := flag.Arg(1)
	channel := flag.Arg(2)
	s3Bucket := flag.Arg(3)
	if flag.NArg() != 4 {
		printUsage()
		os.Exit(1)
	}
	versions := autoupdate.VersionFile{
		Versions:    make([]string, 0),
		LastVersion: "",
	}
	versionFileKey := autoupdate.GetVersionFileKey(appName, channel)

	emptyVersionsFile, _ := json.Marshal(versions)
	fmt.Println("Uploading empty " + versionFileKey + " file")

	err := autoupdate.UploadS3File(awsConfig, s3Bucket, versionFileKey, bytes.NewReader(emptyVersionsFile))
	if err != nil {
		panic(err)
	}
}
func downloadLatest(awsConfig *aws.Config) {
	appName := flag.Arg(1)
	channel := flag.Arg(2)
	s3Bucket := flag.Arg(3)
	outputName := flag.Arg(4)
	if flag.NArg() != 5 {
		printUsage()
		os.Exit(1)
	}
	versionFileKey := autoupdate.GetVersionFileKey(appName, channel)
	fmt.Println("Fetching " + versionFileKey + " file")

	versionFile, err := autoupdate.GetS3File(awsConfig, s3Bucket, versionFileKey, false, nil)

	if err != nil {
		panic(err)
	}

	var versions autoupdate.VersionFile

	err = json.Unmarshal(versionFile, &versions)
	if err != nil {
		panic(err)
	}

	releaseFileKey := autoupdate.GetFileKey(appName, channel, versions.LastVersion+"/"+"windows-amd64")
	fmt.Println("Fetching " + releaseFileKey + " file")
	releaseFile, err := autoupdate.GetS3File(awsConfig, s3Bucket, releaseFileKey, false, func(i int) {
		fmt.Printf("\rDownloading file %v: %v%%", releaseFileKey, i)
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("... Done")

	err = ioutil.WriteFile(outputName, releaseFile, 0755)
	if err != nil {
		panic(err)
	}
}

func release(awsConfig *aws.Config) {
	appName := flag.Arg(1)
	channel := flag.Arg(2)
	s3Bucket := flag.Arg(3)
	releasesDir := flag.Arg(4)
	releaseTag := flag.Arg(5)
	if flag.NArg() != 6 {
		printUsage()
		os.Exit(1)
	}

	files, err := ioutil.ReadDir(releasesDir)

	if err != nil {
		panic(err)
	}

	for _, file := range files {

		fileBody, err := os.Open(releasesDir + "/" + file.Name())
		if err != nil {
			panic(err)
		}
		fileKey := autoupdate.GetFileKey(appName, channel, releaseTag+"/"+file.Name())
		fmt.Println("Uploading " + fileKey + " file")
		err = autoupdate.UploadS3File(
			awsConfig,
			s3Bucket,
			fileKey,
			fileBody,
		)
		if err != nil {
			panic(err)
		}
		err = fileBody.Close()
		if err != nil {
			panic(err)
		}
	}

	versionFileKey := autoupdate.GetVersionFileKey(appName, channel)

	fmt.Println("Getting " + versionFileKey + " file")
	versionFile, _ := autoupdate.GetS3File(awsConfig, s3Bucket, versionFileKey, false, nil)

	var versions autoupdate.VersionFile

	err = json.Unmarshal(versionFile, &versions)
	if err != nil {
		panic(err)
	}

	versions.Versions = append(versions.Versions, releaseTag)
	versions.LastVersion = releaseTag

	updatedVersionsFile, _ := json.Marshal(versions)

	fmt.Println("Uploading " + versionFileKey + " file")
	err = autoupdate.UploadS3File(awsConfig, s3Bucket, versionFileKey, bytes.NewReader(updatedVersionsFile))
	if err != nil {
		panic(err)
	}
}

func main() {
	region := flag.String("region", "us-east-1", "S3 region")
	endpoint := flag.String("endpoint", "", "S3 endpoint URL")
	awsDisableSSL := flag.Bool("disable-ssl", false, "Disable SSL")
	s3ForcePathStyle := flag.Bool("s3-force-path-style", false, "S3 force path style")
	flag.Parse()

	awsConfig := &aws.Config{
		Region: aws.String(*region),
	}
	if *endpoint != "" {
		awsConfig.Endpoint = aws.String(*endpoint)
	}
	if *awsDisableSSL {
		awsConfig.DisableSSL = aws.Bool(true)
	}
	if *s3ForcePathStyle {
		awsConfig.S3ForcePathStyle = aws.Bool(true)
	}

	applicationMode := flag.Arg(0)

	switch applicationMode {
	case "release":
		release(awsConfig)
	case "init":
		initBucket(awsConfig)
	case "download-latest":
		downloadLatest(awsConfig)
	default:
		printUsage()
		os.Exit(1)
	}
}
