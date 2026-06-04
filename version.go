// Copyright (c) 2026 Steve Simenic. All rights reserved.
// Use of this source code is governed by a GNU AGPLv3 license
// that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"
	"strings"
)

const (
	versionFile = "updater.txt"
)

func localVersion() (string, error) {
	ver, err := os.ReadFile(versionFile)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", fmt.Errorf("local version file unreadable: %w", err)
	}
	return strings.TrimSpace(string(ver)), nil
}

func writeVersion(ver string) error {
	data := []byte(ver)
	return os.WriteFile(versionFile, data, 0644)
}
