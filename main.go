// Copyright (c) 2026 Steve Simenic. All rights reserved.
// Use of this source code is governed by a GNU AGPLv3 license
// that can be found in the LICENSE file.

// Package main provides the main execution logic for the Augustus updater,
// a utility to download the latest Augustus unstable build and run it.
package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"golang.org/x/mod/semver"
)

const (
	updaterURL = "https://github.com/orffen/augustus-updater"
)

var (
	version = "0.0.0-dev"
)

func main() {
	fmt.Println("Augustus Updater", version)
	updaterVer, err := latestUpdaterVersion()
	if err != nil {
		showError("Couldn't check updater version:", err)
	}
	if semver.Compare(updaterVer, version) == 1 {
		fmt.Println("New updater version", updaterVer, "is available! Download from", updaterURL)
		time.Sleep(3 * time.Second)
	}
	fmt.Println("Checking for latest Augustus unstable version...")
	lastURL, err := localVersion()
	if err != nil {
		showError("Couldn't read local version:", err)
	}
	if url, err := getDownloadURL(regex); err != nil {
		showError("Couldn't find download URL:", err)
		if lastURL == "" {
			fatalError("No current version installed and unable to download latest.")
		}
	} else {
		if lastURL != url {
			fmt.Println("Updating, please wait...")
			if err := applyUpdate(url); err != nil {
				fatalError("Couldn't apply update:", err)
			}
			if err := writeVersion(url); err != nil {
				showError("Couldn't write local version:", err)
			}
		} else {
			fmt.Println("Already up-to-date.")
		}
	}
	runProgram()
}

func showError(v ...any) {
	args := append([]any{"Warning:"}, v...)
	fmt.Fprintln(os.Stderr, args...)
}

func fatalError(v ...any) {
	args := append([]any{"Fatal Error:"}, v...)
	_ = writeVersion(fmt.Sprintf("%s\nWill redownload when next run.", args)) // force a re-download next run
	fmt.Fprintln(os.Stderr, args...)
	fmt.Println("Press ENTER key to quit...")
	reader := bufio.NewScanner(os.Stdin)
	_ = reader.Scan()
	os.Exit(1)
}
