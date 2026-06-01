// Copyright (c) 2026 Steve Simenic. All rights reserved.
// Use of this source code is governed by a GNU AGPLv3 license
// that can be found in the LICENSE file.

//go:generate goversioninfo

// Package main provides the main execution logic for the Augustus updater,
// a utility to download the latest Augustus unstable build and run it.
package main

import (
	"fmt"
	"os"
)

func main() {
	var err error
	var lastURL string
	var url string

	lastURL, err = localVersion()
	if err != nil {
		fmt.Println("Couldn't read local version:", err)
	}
	url, err = getDownloadURL(regex)
	if err != nil {
		fmt.Println("Error:", err)
		if lastURL == "" {
			fmt.Println("Fatal Error: No current version installed and unable to download latest.")
			os.Exit(1)
		}
	} else {
		if lastURL != url {
			err = applyUpdate(url)
			if err != nil {
				fmt.Println("Fatal Error applying update:", err)
				os.Exit(1)
			}
			err = writeVersion(url)
			if err != nil {
				fmt.Println("Couldn't write local version:", err)
			}
		}
	}
	runProgram()
}
