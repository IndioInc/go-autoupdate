package autoupdate

import (
	"io/ioutil"
	"os"
	"strings"
)

func getVersionFilePath(updater *Updater) (string, error) {
	dir, err := getExecutableDirectory()
	if err != nil {
		return "", err
	}

	return dir + "/" + updater.versionFilePath, nil
}

func wasUpdated(updater *Updater, latestVersion string) (bool, error) {
	versionFilePath, err := getVersionFilePath(updater)
	if err != nil {
		return false, err
	}

	f, err := ioutil.ReadFile(versionFilePath)
	if os.IsNotExist(err) {
		err = updateCurrentVersion(updater, latestVersion)
		return false, err
	} else if err != nil {
		return false, err
	}

	currentVersion := strings.TrimSpace(string(f))

	return currentVersion != latestVersion, err
}

func updateCurrentVersion(updater *Updater, release string) error {
	versionFilePath, err := getVersionFilePath(updater)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(versionFilePath, []byte(release), 0644)
}
