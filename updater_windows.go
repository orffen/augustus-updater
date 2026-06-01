// Copyright (c) 2026 Steve Simenic. All rights reserved.
// Use of this source code is governed by a GNU AGPLv3 license
// that can be found in the LICENSE file.

package main

import (
	"archive/zip"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

const (
	exe     = "./augustus.exe"
	regex   = `href="([^"]+windows\.zip)"`
	outFile = "temp.zip"
)

func applyUpdate(url string) error {
	var err error
	replacer := strings.NewReplacer(
		"augustus", "assets",
		"windows", "development",
	)
	assetsURL := replacer.Replace(url)
	defer func() { _ = os.Remove(outFile) }()
	for _, u := range []string{url, assetsURL} {
		if err = downloadUpdate(u, outFile); err != nil {
			return err
		}
		if err = unzip(outFile); err != nil {
			return err
		}
		if err = os.Remove(outFile); err != nil {
			return fmt.Errorf("couldn't remove temporary file %v: %w", outFile, err)
		}
	}
	return nil
}

func runProgram() {
	cmd := exec.Command(exe)
	cmd.Dir = "."
	// CREATE_NEW_PROCESS_GROUP = 0x00000200
	// DETACHED_PROCESS = 0x00000008
	cmd.SysProcAttr = &syscall.SysProcAttr{
		CreationFlags: 0x00000200 | 0x00000008,
	}
	_ = cmd.Start()
	time.Sleep(100 * time.Millisecond)
	os.Exit(0)
}

func unzip(file string) error {
	r, err := zip.OpenReader(file)
	if err != nil {
		return fmt.Errorf("couldn't open zip %v: %w", file, err)
	}
	defer func() { _ = r.Close() }()
	if err := os.CopyFS(".", r); err != nil {
		return fmt.Errorf("couldn't extract zip %v: %w", file, err)
	}
	return nil
}
