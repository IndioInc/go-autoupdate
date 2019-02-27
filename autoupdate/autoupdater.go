package autoupdate

import (
	"os"
	"time"
)

// Updater configuration. Use NewUpdater() to construct this, then use setters to change some of the additional params
type updater struct {
	s3Bucket          string
	channel           string
	appName           string
	checkInterval     int
	releasesDirectory string
	versionFilePath   string
}

func NewUpdater(s3Bucket string, channel string, appName string, versionFile string) *updater {
	return &updater{
		s3Bucket:          s3Bucket,
		channel:           channel,
		appName:           appName,
		checkInterval:     10,
		releasesDirectory: "releases",
		versionFilePath:   versionFile,
	}
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
func RunAutoupdater(updater *updater, shutdownCallback func(error)) {
	err := func() error {
		for {
			changed, err := IsNewVersionAvailable(updater)
			if err != nil {
				return err
			}

			if changed {
				err := UpdateApplication(updater)
				if err != nil {
					return err
				}
				return nil
			}
			time.Sleep(time.Duration(updater.checkInterval) * time.Second)
		}
	}()
	shutdownCallback(err)

	exitCode := 0
	if err != nil {
		exitCode = 1
	}

	os.Exit(exitCode)
}

/*
Runs update check once. Returns if there's new version available to be downloaded. If so, you can run UpdateApplication function
*/
func IsNewVersionAvailable(updater *updater) (bool, error) {
	latestTag, err := getLatestVersionTag(updater)
	if err != nil {
		return false, err
	}
	return wasUpdated(updater, latestTag)
}

/*
Updates the application. To see if you should run this, first you should call `IsNewersionAvailable`
*/
func UpdateApplication(updater *updater) error {
	releaseVersion, err := downloadRelease(updater)
	if err != nil {
		return err
	}
	err = swapReleaseFiles()
	if err != nil {
		return err
	}

	err = updateCurrentVersion(updater, releaseVersion)
	if err != nil {
		return err
	}
	return nil
}
