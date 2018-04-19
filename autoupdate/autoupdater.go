package autoupdate

import (
	"time"
)

type Updater struct {
	S3Bucket                string
	Channel                 string
	AppName                 string
	CheckInterval           int
	ReleasesDirectory       string
	UnauthenticatedDownload bool
}

var cmd *command

func startApplication(filename string) {
	StopRunningApplication()

	cmd = createCommand(filename)

	err := cmd.Start()
	if err != nil {
		panic(err)
	}
	cmd.listenForStop()
}

func RunAutoupdater(updater Updater) {
	for {
		if hasS3FileChanged(updater) {
			releaseFilename := downloadLatestRelease(updater)

			startApplication(releaseFilename)
		}
		time.Sleep(time.Duration(updater.CheckInterval) * time.Second)
	}
}

func StopRunningApplication() {
	if cmd != nil {
		cmd.stop()
	}
}
