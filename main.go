// Copyright (c) 2026 Steve Simenic. All rights reserved.
// Use of this source code is governed by a GNU AGPLv3 license
// that can be found in the LICENSE file.

// Package main provides the main execution logic for the Augustus updater,
// a utility to download the latest Augustus unstable build and run it.
package main

import (
	"fmt"
)

func main() {
	lastURL, err := localVersion()
	if err != nil {
		fmt.Println("Couldn't read local version:", err)
	}
	if url, err := getDownloadURL(regex); err != nil {
		fmt.Println("Error:", err)
		if lastURL == "" {
			fatalError("Fatal Error: No current version installed and unable to download latest.")
		}
	} else {
		if lastURL != url {
			if err := applyUpdate(url); err != nil {
				fatalError("Fatal Error applying update:", err)
			}
			if err := writeVersion(url); err != nil {
				fmt.Println("Couldn't write local version:", err)
			}
		}
	}
	runProgram()
}
