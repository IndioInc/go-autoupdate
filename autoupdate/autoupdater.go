package autoupdate

import (
	"os"
	"os/exec"
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

var cmd *exec.Cmd

func startApplication(filename string) {
	if cmd != nil {
		cmd.Process.Kill()
		cmd.Process.Wait()
	}
	cmd = exec.Command(filename)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Start()
	if err != nil {
		panic(err)
	}

	go func() {
		cmd.Wait()
		if cmd.ProcessState != nil && cmd.ProcessState.Exited() {
			if cmd.ProcessState.Success() {
				os.Exit(0)
			}
		}
	}()
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
