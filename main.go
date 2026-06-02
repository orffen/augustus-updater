// Copyright (c) 2026 Steve Simenic. All rights reserved.
// Use of this source code is governed by a GNU AGPLv3 license
// that can be found in the LICENSE file.

// Package main provides the main execution logic for the Augustus updater,
// a utility to download the latest Augustus unstable build and run it.
package main

import "fmt"

func main() {
	fmt.Println("Augustus Updater checking for latest unstable version...")
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
			if err := applyUpdate(url); err != nil {
				fatalError("Couldn't apply update:", err)
			}
			if err := writeVersion(url); err != nil {
				showError("Couldn't write local version:", err)
			}
		}
	}
	runProgram()
}
