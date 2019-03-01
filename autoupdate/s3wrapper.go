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

func getS3Downloader(anonymousSession bool) (*s3manager.Downloader, error) {
	sess, err := getAwsSession(anonymousSession)
	if err != nil {
		return nil, err
	}

	return s3manager.NewDownloader(sess), nil
}

func getS3Uploader() (*s3manager.Uploader, error) {
	sess, err := getAwsSession(false)
	if err != nil {
		return nil, err
	}

	return s3manager.NewUploader(sess), nil
}

func getFileSize(s3Bucket string, key string, anonymousSession bool) (int64, error) {
	awsSession, err := getAwsSession(anonymousSession)
	if err != nil {
		return 0, err
	}
	svc := s3.New(awsSession)
	resp, err := svc.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(s3Bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return 0, err
	}

	return *resp.ContentLength, nil
}

// Gets an S3 file and returns the body and error.
func GetS3File(s3Bucket string, key string, anonymousSession bool, progressCallback func(int)) ([]byte, error) {
	s3Downloader, err := getS3Downloader(anonymousSession)
	if err != nil {
		return nil, err
	}
	temp := aws.NewWriteAtBuffer([]byte{})
	size, err := getFileSize(s3Bucket, key, anonymousSession)
	if err != nil {
		return nil, err
	}

	writer := &progressWriter{writer: temp, size: size, written: 0, progressCallback: progressCallback}

	params := &s3.GetObjectInput{
		Bucket: aws.String(s3Bucket),
		Key:    aws.String(key),
	}
	if _, err := s3Downloader.Download(writer, params); err != nil {
		return nil, err
	}

	return writer.writer.Bytes(), nil
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

func getLatestVersionTag(updater *Updater) (string, error) {
	var versions VersionFile

	file, err := GetS3File(updater.s3Bucket, GetVersionFileKey(updater.appName, updater.channel), true, nil)
	if err != nil {
		return "", err
	}
	err = json.Unmarshal(file, &versions)

	return versions.LastVersion, err
}

func downloadRelease(updater *Updater, progressCallback func(int)) (string, error) {
	version, err := getLatestVersionTag(updater)
	if err != nil {
		return "", err
	}

	fileKey := getReleaseFileKey(updater.appName, updater.channel, version)

	file, err := GetS3File(updater.s3Bucket, fileKey, true, progressCallback)
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

	_, err = outFile.Write(file)

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
