package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/IndioInc/go-autoupdate/autoupdate"
	"io/ioutil"
	"os"
)

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("\tgo-autoupdate help")
	fmt.Println("\tgo-autoupdate init <appName> <channel> <s3Bucket>")
	fmt.Println("\tgo-autoupdate release <appName> <channel> <s3Bucket> <releasesDir> <releasesTag>")
	fmt.Println("\tgo-autoupdate download-latest <appName> <channel> <s3Bucket> <outputName>")
}

func initBucket() {
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

	err := autoupdate.UploadS3File(s3Bucket, versionFileKey, bytes.NewReader(emptyVersionsFile))
	if err != nil {
		panic(err)
	}
}
func downloadLatest() {
	appName := flag.Arg(1)
	channel := flag.Arg(2)
	s3Bucket := flag.Arg(3)
	outputName := flag.Arg(4)
	if flag.NArg() != 5 {
		printUsage()
		os.Exit(1)
	}
	versionFileKey := autoupdate.GetVersionFileKey(appName, channel)
	fmt.Println("Fetching " + versionFileKey + "file")

	versionFile, err := autoupdate.GetS3File(s3Bucket, versionFileKey, false, nil)

	if err != nil {
		panic(err)
	}

	var versions autoupdate.VersionFile

	err = json.Unmarshal(versionFile, &versions)
	if err != nil {
		panic(err)
	}

	releaseFileKey := autoupdate.GetFileKey(appName, channel, versions.LastVersion+"/"+"windows-amd64")
	fmt.Println("Fetching " + releaseFileKey + "file")
	releaseFile, err := autoupdate.GetS3File(s3Bucket, releaseFileKey, false, func(i int) {
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

func release() {
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
		err = autoupdate.UploadS3File(s3Bucket,
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
	versionFile, _ := autoupdate.GetS3File(s3Bucket, versionFileKey, false, nil)

	var versions autoupdate.VersionFile

	err = json.Unmarshal(versionFile, &versions)
	if err != nil {
		panic(err)
	}

	versions.Versions = append(versions.Versions, releaseTag)
	versions.LastVersion = releaseTag

	updatedVersionsFile, _ := json.Marshal(versions)

	fmt.Println("Uploading " + versionFileKey + " file")
	err = autoupdate.UploadS3File(s3Bucket, versionFileKey, bytes.NewReader(updatedVersionsFile))
	if err != nil {
		panic(err)
	}
}

func main() {
	flag.Parse()

	applicationMode := flag.Arg(0)

	switch applicationMode {
	case "release":
		release()
	case "init":
		initBucket()
	case "download-latest":
		downloadLatest()
	default:
		printUsage()
		os.Exit(1)
	}
}
