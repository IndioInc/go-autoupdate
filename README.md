# go-autoupdate

## Description 

This package provides way to create auto-updating go-lang applications.

You can achieve that by uploading your binary files into an S3 bucket and listening for changes in a goroutine

## Installation

```
go get -u github.com/IndioInc/go-autoupdate
```

## Usage

### Application Setup

updater/main.go
```go
package main

import (
    "github.com/IndioInc/go-autoupdate/autoupdate"
)

var updater = autoupdate.NewUpdater(
    "company-releases-bucket", 
    "your-app",
    "stable",
    ".app-version",
    "us-east-1",
    nil,
    false,
    false,
)

func main() {
    go autoupdate.RunAutoupdater(updater, func(err error) {
        // gracefully handle shutdown
        if err != nil {
            panic(err)
        }
    })
    
    // do actual application stuff
}
```

This file is a boilerplate. Your actual application code goes to `app/main.go` (for example).
You don't need any boilerplate code in your application. It will be run as a subprocess of the autoupdater.
Stdout and Stderr is being forwarded to the updater process.

### Making releases

First, run `go-autoupdate init your-app stable company-releases-bucket`

Compile your application for all environments you wish the application to work storing them in `releases/` directory (or whatever you set in config). 

Names of the binaries need to follow naming convention of `{{GOOS}}-{{GOARCH}}`.

After you've compiled all the binaries, use `release` command provided by this package to send the release to S3.

`go-autoupdate release your-app stable company-releases-bucket releases $COMMIT_ID`.

All applications listening for changes will see that `version.json` file has changed in the bucket, download the new release for the correct GOOS-GOARCH pair, stop the old process with `SIGTERM` and start the new application version.



