package autoupdate

import (
	"fmt"
	"os"
	"runtime"
)

func getLocalReleaseFilename() (string, error) {
	return os.Executable()
}

func getNewReleaseFilename() (string, error) {
	fileName, err := getLocalReleaseFilename()
	if err != nil {
		return "", err
	}
	return fileName + ".new", nil
}

func getOldReleaseFilename() (string, error) {
	fileName, err := getLocalReleaseFilename()
	if err != nil {
		return "", err
	}
	return fileName + ".old", nil
}

func GetFileKey(appName string, channel string, filename string) string {
	return fmt.Sprintf("%s/%s/%s", appName, channel, filename)
}

func GetVersionFileKey(appName string, channel string) string {
	return GetFileKey(appName, channel, "versions.json")
}

func getOsArch() string {
	return fmt.Sprintf("%s-%s", runtime.GOOS, runtime.GOARCH)
}

func getReleaseFileKey(appName string, channel string, version string) string {
	osArch := getOsArch()
	return GetFileKey(appName, channel, fmt.Sprintf("%s/%s", version, osArch))
}
