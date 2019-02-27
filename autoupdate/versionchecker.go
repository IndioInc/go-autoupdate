package autoupdate

import (
	"io/ioutil"
	"os"
)

func wasUpdated(updater *updater, latestVersion string) (bool, error) {
	f, err := ioutil.ReadFile(updater.versionFilePath)
	if os.IsNotExist(err) {
		err = updateCurrentVersion(updater, latestVersion)
		return false, nil
	} else if err != nil {
		return false, err
	}

	currentVersion := string(f)

	return currentVersion != latestVersion, err
}

func updateCurrentVersion(updater *updater, release string) error {
	return ioutil.WriteFile(updater.versionFilePath, []byte(release), 0644)
}
