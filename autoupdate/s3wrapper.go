package autoupdate

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"io"
	"log"
	"os"
)

type VersionFile struct {
	Versions    []string `json:"versions"`
	LastVersion string   `json:"lastVersion"`
}

func getAwsSession(anonymousSession bool) *session.Session {
	var awsSession *session.Session
	config := &aws.Config{Region: aws.String("us-east-1")}
	if anonymousSession {
		config.Credentials = credentials.AnonymousCredentials
	}

	awsSession, err := session.NewSession(config)
	if err != nil {
		log.Fatal(err)
	}

	return awsSession
}

func getS3Client(unauthenticatedDownload bool) *s3.S3 {
	sess := getAwsSession(unauthenticatedDownload)

	return s3.New(sess)
}

func getS3Uploader() *s3manager.Uploader {
	sess := getAwsSession(false)

	return s3manager.NewUploader(sess)
}

var lastETag = ""

func hasS3FileChanged(updater Updater) bool {
	s3Client := getS3Client(updater.UnauthenticatedDownload)

	result, err := s3Client.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(updater.S3Bucket),
		Key:    aws.String(GetVersionFileKey(updater.AppName, updater.Channel)),
	})

	checkError(err)

	return *result.ETag != lastETag
}

// Gets an S3 file and returns the body and ETag
func GetS3File(s3Bucket string, key string, unauthenticatedDownload bool) (io.ReadCloser, string) {
	s3Client := getS3Client(unauthenticatedDownload)

	result, err := s3Client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(s3Bucket),
		Key:    aws.String(key),
	})

	checkError(err)

	file := result.Body

	return file, *result.ETag
}

func UploadS3File(s3Bucket string, key string, file io.Reader) {
	uploader := getS3Uploader()

	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(s3Bucket),
		Key:    aws.String(key),
		Body:   file,
	})

	if err != nil {
		panic(err)
	}
}

func getLatestVersionTag(updater Updater) string {
	var versions VersionFile

	file, eTag := GetS3File(updater.S3Bucket, GetVersionFileKey(updater.AppName, updater.Channel), updater.UnauthenticatedDownload)

	lastETag = eTag

	err := json.NewDecoder(file).Decode(&versions)

	checkError(err)

	return versions.LastVersion
}

// Downloads latest release and returns the filename
func downloadLatestRelease(updater Updater) string {
	version := getLatestVersionTag(updater)

	fileKey := getReleaseFileKey(updater.AppName, updater.Channel, version)

	file, _ := GetS3File(updater.S3Bucket, fileKey, updater.UnauthenticatedDownload)

	ensureDirectoryExists(updater.ReleasesDirectory)
	releaseFilename := getLocalReleaseFilename(updater.ReleasesDirectory, version)
	outFile, err := os.OpenFile(releaseFilename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)

	checkError(err)
	defer outFile.Close()

	_, err = io.Copy(outFile, file)
	checkError(err)

	return releaseFilename
}
