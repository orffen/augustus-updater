// Copyright (c) 2026 Steve Simenic. All rights reserved.
// Use of this source code is governed by a GNU AGPLv3 license
// that can be found in the LICENSE file.

package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
)

const (
	exe     = "Augustus Unstable.AppImage"
	outFile = "Augustus Unstable.tmp"
	regex   = `href="([^"]+linux\.AppImage)"`
)

func applyUpdate(url string) error {
	defer func() { _ = os.Remove(outFile) }()
	if err := downloadUpdate(context.Background(), url, outFile); err != nil {
		return err
	}
	if err := os.Chmod(outFile, 0755); err != nil {
		_ = os.Remove(outFile)
		return fmt.Errorf("couldn't make AppImage executable: %w", err)
	}
	if err := os.Rename(outFile, exe); err != nil {
		_ = os.Remove(outFile)
		return fmt.Errorf("couldn't deploy new AppImage: %w", err)
	}
	return nil
}

func runProgram() error {
	cmd := exec.Command("./" + exe)
	if err := cmd.Start(); err != nil {
		return err
	}
	return nil
}
