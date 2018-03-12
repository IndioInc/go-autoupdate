package autoupdate

import (
	"fmt"
	"os"
	"runtime"
)

func getLocalReleaseFilename(releasesDirectory string, version string) string {
	fileSuffix := ""
	if runtime.GOOS == "windows" {
		fileSuffix = ".exe"
	}
	return fmt.Sprintf("%s/%s%s", releasesDirectory, version, fileSuffix)
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

func ensureDirectoryExists(directoryName string) {
	os.MkdirAll(directoryName, 0775)
}
