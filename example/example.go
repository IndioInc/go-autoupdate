package main

import (
	"fmt"
	"github.com/IndioInc/go-autoupdate/autoupdate"
	"runtime"
	"time"
)

var updater = autoupdate.NewUpdater(
	"company-releases-bucket",
	"stable",
	"your-app",
	".example-version",
)

func main() {
	fmt.Println("Starting application Application")
	fmt.Println(runtime.GOOS, runtime.GOARCH)
	go autoupdate.RunAutoupdater(updater, func(err error) {
		// gracefully handle shutdown
		if err != nil {
			panic(err)
		}
	})
	for {
		time.Sleep(1 * time.Second)
		fmt.Println("Barr")
	}
}
