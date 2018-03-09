package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/IndioInc/go-autoupdate"
	"io/ioutil"
	"os"
)

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("\trelease help")
	fmt.Println("\trelease add <appName> <channel> <s3Bucket> <releasesDir> <releasesTag>")
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

	versionFileKey := autoupdate.GetVersionFileKey(appName, channel)

	fmt.Println("Getting " + versionFileKey + " file")
	versionFile, _ := autoupdate.GetS3File(s3Bucket, versionFileKey)

	var versions autoupdate.VersionFile

	json.NewDecoder(versionFile).Decode(&versions)

	versions.Versions = append(versions.Versions, releaseTag)
	versions.LastVersion = releaseTag

	updatedVersionsFile, _ := json.Marshal(versions)

	fmt.Println("Uploading " + versionFileKey + " file")
	autoupdate.UploadS3File(s3Bucket, versionFileKey, bytes.NewReader(updatedVersionsFile))

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
		autoupdate.UploadS3File(s3Bucket,
			fileKey,
			fileBody,
		)
		fileBody.Close()
	}
}

func main() {
	flag.Parse()

	applicationMode := flag.Arg(0)

	switch applicationMode {
	case "add":
		release()
	default:
		printUsage()
		os.Exit(1)
	}
}
