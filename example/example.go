package main

import (
	"fmt"
	"github.com/IndioInc/go-autoupdate/autoupdate"
	"os"
	"runtime"
	"time"
)

var updater = autoupdate.NewUpdater("company-releases-bucket", "stable", "your-app")

func main() {
	fmt.Println("Starting application Application")
	fmt.Println(runtime.GOOS, runtime.GOARCH)
	go func() {
		err := autoupdate.RunAutoupdater(updater)
		if err != nil {
			panic(err)
		}
		// gracefully handle shutdown
		os.Exit(0)
	}()
	for {
		time.Sleep(1 * time.Second)
		fmt.Println("Barr")
	}
}
