package autoupdate

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"io"
	"os"
)

type VersionFile struct {
	Versions    []string `json:"versions"`
	LastVersion string   `json:"lastVersion"`
}

func getAwsSession(anonymousSession bool) (*session.Session, error) {
	config := &aws.Config{Region: aws.String("us-east-1")}
	if anonymousSession {
		config.Credentials = credentials.AnonymousCredentials
	}

	return session.NewSession(config)
}

func getS3Client(anonymousSession bool) (*s3.S3, error) {
	sess, err := getAwsSession(anonymousSession)
	if err != nil {
		return nil, err
	}

	return s3.New(sess), nil
}

func getS3Uploader() (*s3manager.Uploader, error) {
	sess, err := getAwsSession(false)
	if err != nil {
		return nil, err
	}

	return s3manager.NewUploader(sess), nil
}

// Gets an S3 file and returns the body and ETag
func GetS3File(s3Bucket string, key string, anonymousSession bool) (io.ReadCloser, error) {
	s3Client, err := getS3Client(anonymousSession)
	if err != nil {
		return nil, err
	}
	fmt.Println(s3Bucket, key)
	result, err := s3Client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(s3Bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		return nil, err
	}

	file := result.Body

	return file, nil
}

func UploadS3File(s3Bucket string, key string, file io.Reader) error {
	uploader, err := getS3Uploader()
	if err != nil {
		return err
	}

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(s3Bucket),
		Key:    aws.String(key),
		Body:   file,
	})

	return err
}

func getLatestVersionTag(updater *updater) (string, error) {
	var versions VersionFile

	file, err := GetS3File(updater.s3Bucket, GetVersionFileKey(updater.appName, updater.channel), true)
	if err != nil {
		return "", err
	}

	err = json.NewDecoder(file).Decode(&versions)

	return versions.LastVersion, err
}

func downloadRelease(updater *updater) (string, error) {
	version, err := getLatestVersionTag(updater)
	if err != nil {
		return "", err
	}

	fileKey := getReleaseFileKey(updater.appName, updater.channel, version)

	file, err := GetS3File(updater.s3Bucket, fileKey, true)
	if err != nil {
		return "", err
	}

	releaseFilename, err := getNewReleaseFilename()
	if err != nil {
		return "", err
	}
	fmt.Println(releaseFilename)
	outFile, err := os.OpenFile(releaseFilename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return "", err
	}

	fmt.Println("saving file", releaseFilename)

	defer outFile.Close()

	_, err = io.Copy(outFile, file)

	return version, err
}

func swapReleaseFiles() error {
	oldFileName, err := getOldReleaseFilename()
	if err != nil {
		return err
	}
	releaseFilename, err := getLocalReleaseFilename()
	if err != nil {
		return err
	}
	newFileName, err := getNewReleaseFilename()
	if err != nil {
		return err
	}
	err = os.Rename(releaseFilename, oldFileName)
	if err != nil {
		return err
	}
	err = os.Rename(newFileName, releaseFilename)

	return err
}
