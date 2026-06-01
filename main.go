// Copyright (c) 2026 Steve Simenic. All rights reserved.
// Use of this source code is governed by a GNU AGPLv3 license
// that can be found in the LICENSE file.

// Package main provides the main execution logic for the Augustus updater,
// a utility to download the latest Augustus unstable build and run it.
package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
)

const (
	updateURL   = "https://josecadete.net/"
	versionFile = "download_url.txt"
)

func getDownloadURL() (string, error) {
	resp, err := http.Get(updateURL)
	if err != nil {
		return "", fmt.Errorf("failed to fetch %v: %w", updateURL, err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("bad server response from %v: %w", updateURL, err)
	}
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("couldn't read page body from %v: %w", updateURL, err)
	}
	htmlContent := string(bodyBytes)
	rex := regexp.MustCompile(regex)
	matches := rex.FindStringSubmatch(htmlContent)
	if len(matches) < 1 {
		return "", fmt.Errorf("couldn't find filename to match regex %v", regex)
	}
	return updateURL + matches[1], nil
}

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
