package autoupdate

import (
	"os"
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
}

var cmd *exec.Cmd

func startApplication(filename string) {
	if cmd != nil {
		cmd.Process.Signal(syscall.SIGTERM)
		_, err := cmd.Process.Wait()
		checkError(err)
	}
	cmd = exec.Command(filename)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Start()
	if err != nil {
		panic(err)
	}
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
