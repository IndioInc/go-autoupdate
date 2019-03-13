package main

import (
	"fmt"
	"time"

	"github.com/IndioInc/go-autoupdate/autoupdate"
)

var updater = autoupdate.NewUpdater(
	"company-releases-bucket",
	"stable",
	"your-app",
	".example-version",
)

func main() {
	fmt.Println("Starting Application")
	go autoupdate.RunAutoupdater(updater, func(err error) {
		// gracefully handle shutdown
		if err != nil {
			panic(err)
		}
	})
	for {
		time.Sleep(1 * time.Second)
	}
}
