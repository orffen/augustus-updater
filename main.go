// Copyright (c) 2026 Steve Simenic. All rights reserved.
// Use of this source code is governed by a GNU AGPLv3 license
// that can be found in the LICENSE file.

// Package main provides the main execution logic for the Augustus updater,
// a utility to download the latest Augustus unstable build and run it.
package main

import (
	"fmt"
	"os"
	"strings"
)

const (
	updateURL   = "https://josecadete.net/"
	versionFile = "download_url.txt"
)

func localVersion() (string, error) {
	ver, err := os.ReadFile(versionFile)
	if err != nil {
		return "-1", fmt.Errorf("failed to read %v: %w", versionFile, err)
	}
	return strings.TrimSpace(string(ver)), nil
}

func writeVersion(ver string) error {
	data := []byte(ver)
	return os.WriteFile(versionFile, data, 0644)
}

func main() {
	var err error
	var lastURL string
	var url string

	lastURL, err = localVersion()
	if err != nil {
		fmt.Println("Couldn't read local version:", err)
	}
	url, err = getDownloadURL()
	if err != nil {
		fmt.Println("Error:", err)
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
