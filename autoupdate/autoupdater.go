package autoupdate

import (
	"fmt"
	"os"
	"time"
)

type updater struct {
	s3Bucket          string
	channel           string
	appName           string
	checkInterval     int
	releasesDirectory string
}

func NewUpdater(s3Bucket string, channel string, appName string) *updater {
	return &updater{s3Bucket: s3Bucket, channel: channel, appName: appName, checkInterval: 10, releasesDirectory: "releases"}
}

func (u *updater) SetInterval(interval int) {
	u.checkInterval = interval
}

func (u *updater) SetReleaseDirectory(releaseDirectory string) {
	u.releasesDirectory = releaseDirectory
}

/*
Starts autoupdater. When release file has changed, the application gets downloaded and then stopped.
It is developer's job to make sure the application gets restarted (most of the time using a service)
*/
func RunAutoupdater(updater *updater) error {
	for {
		latestTag, err := getLatestVersionTag(updater)
		if err != nil {
			return err
		}
		changed, err := wasUpdated(latestTag)
		if err != nil {
			return err
		}
		fmt.Println(changed)
		if changed {
			err := downloadRelease(updater, latestTag)
			if err != nil {
				return err
			}
			err = swapReleaseFiles()
			if err != nil {
				return err
			}

			err = updateCurrentVersion(latestTag)
			if err != nil {
				return err
			}
			os.Exit(0)
		}
		time.Sleep(time.Duration(updater.checkInterval) * time.Second)
	}
}
