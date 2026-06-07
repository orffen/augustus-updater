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
	"strings"
	"time"

	"golang.org/x/mod/semver"
)

const (
	colorReset   = "\033[0m"  // ANSI reset
	colorError   = "\033[31m" // ANSI red
	colorUpdate  = "\033[32m" // ANSI green
	colorWarning = "\033[33m" // ANSI yellow
	updaterURL   = "https://github.com/orffen/augustus-updater"
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
		fmt.Println(colorUpdate+"New updater version", updaterVer, "is available! Download from", updaterURL+colorReset)
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
			fmt.Println("Already up to date.")
		}
	}
	if err := runProgram(); err != nil {
		fatalError("Couldn't start Augustus Unstable:", err)
	}
	time.Sleep(100 * time.Millisecond)
	os.Exit(0)
}

func showError(v ...any) {
	args := strings.TrimSpace(fmt.Sprintln(v...))
	fmt.Fprint(os.Stderr, colorWarning)
	fmt.Fprintln(os.Stderr, "Warning:", args)
	fmt.Fprint(os.Stderr, colorReset)
}

func fatalError(v ...any) {
	args := strings.TrimSpace(fmt.Sprintln(v...))
	_ = writeVersion(fmt.Sprintf("Fatal Error: %s\nWill redownload when next run.", args)) // force a redownload next run
	fmt.Fprint(os.Stderr, colorError)
	fmt.Fprintln(os.Stderr, "Fatal Error:", args)
	fmt.Fprint(os.Stderr, colorReset)
	fmt.Println("Press ENTER key to quit...")
	reader := bufio.NewScanner(os.Stdin)
	_ = reader.Scan()
	os.Exit(1)
}
