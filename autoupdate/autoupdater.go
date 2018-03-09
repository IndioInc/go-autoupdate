package autoupdate

import (
	"os/exec"
	"syscall"
	"time"
)

type Updater struct {
	S3Bucket          string
	Channel           string
	AppName           string
	CheckInterval     int
	ReleasesDirectory string
	CurrentVersion    string
}

var cmd *exec.Cmd

func startApplication(filename string) {
	if cmd != nil {
		cmd.Process.Signal(syscall.SIGTERM)
		_, err := cmd.Process.Wait()
		checkError(err)
	}
	cmd = exec.Command(filename)

	cmd.Start()
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
