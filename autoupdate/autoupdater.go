package autoupdate

import (
	"os"
	"time"
)

// Updater configuration. Use NewUpdater() to construct this, then use setters to change some of the additional params
type Updater struct {
	s3Bucket          string
	channel           string
	appName           string
	checkInterval     int
	versionFilePath   string
}

func NewUpdater(s3Bucket string, channel string, appName string, versionFile string) *Updater {
	return &Updater{
		s3Bucket:        s3Bucket,
		channel:         channel,
		appName:         appName,
		checkInterval:   10,
		versionFilePath: versionFile,
	}
}

func (u *Updater) SetInterval(interval int) {
	u.checkInterval = interval
}

/*
Starts autoupdater. When release file has changed, the application gets downloaded and then stopped.
It is developer's job to make sure the application gets restarted (most of the time using a service)
*/
func RunAutoupdater(updater *Updater, shutdownCallback func(error)) {
	err := func() error {
		for {
			changed, err := IsNewVersionAvailable(updater)
			if err != nil {
				return err
			}

			if changed {
				err := UpdateApplication(updater, nil)
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
func IsNewVersionAvailable(updater *Updater) (bool, error) {
	latestTag, err := getLatestVersionTag(updater)
	if err != nil {
		return false, err
	}
	return wasUpdated(updater, latestTag)
}

/*
Updates the application. To see if you should run this, first you should call `IsNewersionAvailable`
*/
func UpdateApplication(updater *Updater, progressCallback func(int)) error {
	releaseVersion, err := downloadRelease(updater, progressCallback)
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
