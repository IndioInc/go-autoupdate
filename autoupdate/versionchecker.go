package autoupdate

import (
	"io/ioutil"
	"os"
)

const versionFileName = ".go-autoupdate-version"

func wasUpdated(latestVersion string) (bool, error) {
	f, err := ioutil.ReadFile(versionFileName)
	if os.IsNotExist(err) {
		err = updateCurrentVersion(latestVersion)
		return false, nil
	} else if err != nil {
		return false, err
	}

	currentVersion := string(f)

	return currentVersion != latestVersion, err
}

func updateCurrentVersion(release string) error {
	return ioutil.WriteFile(versionFileName, []byte(release), 0644)
}
