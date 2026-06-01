// Copyright (c) 2026 Steve Simenic. All rights reserved.
// Use of this source code is governed by a GNU AGPLv3 license
// that can be found in the LICENSE file.

// Package main implements the platform-specific update process
// for Windows.
package main

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

const (
	exe     = "augustus.exe"
	regex   = `href="([^"]+windows\.zip)"`
	zipFile = "temp.zip"
)

func downloadUpdate(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("download failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad server response from %s: %s", url, resp.Status)
	}
	out, err := os.Create(zipFile)
	if err != nil {
		return fmt.Errorf("couldn't create temp file %v: %w", zipFile, err)
	}
	defer func() { _ = out.Close() }()
	if _, err := io.Copy(out, resp.Body); err != nil {
		return fmt.Errorf("failed to save download: %w", err)
	}
	return nil
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

func applyUpdate(url string) error {
	var err error
	replacer := strings.NewReplacer(
		"augustus", "assets",
		"windows", "development",
	)
	assetsURL := replacer.Replace(url)
	defer func() { _ = os.Remove(zipFile) }()
	for _, u := range []string{url, assetsURL} {
		if err = downloadUpdate(u); err != nil {
			return err
		}
		if err = unzip(zipFile); err != nil {
			return err
		}
		if err = os.Remove(zipFile); err != nil {
			return fmt.Errorf("couldn't remove temporary file %v: %w", zipFile, err)
		}
	}
	return nil
}

func runProgram() {
	cmd := exec.Command(exe)
	// CREATE_NEW_PROCESS_GROUP = 0x00000200
	// DETACHED_PROCESS = 0x00000008
	cmd.SysProcAttr = &syscall.SysProcAttr{
		CreationFlags: 0x00000200 | 0x00000008,
	}
	_ = cmd.Start()
	os.Exit(0)
}
