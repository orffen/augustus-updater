// Copyright (c) 2026 Steve Simenic. All rights reserved.
// Use of this source code is governed by a GNU AGPLv3 license
// that can be found in the LICENSE file.

//go:generate goversioninfo

package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

const (
	exe     = "Augustus Unstable.AppImage"
	outFile = "Augustus Unstable.tmp"
	regex   = `href="([^"]+linux\.AppImage)"`
)

func applyUpdate(url string) error {
	defer func() { _ = os.Remove(outFile) }()
	if err := downloadUpdate(url, outFile); err != nil {
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

func runProgram() {
	cmd := exec.Command("./" + exe)
	_ = cmd.Start()
	time.Sleep(100 * time.Millisecond)
	os.Exit(0)
}
