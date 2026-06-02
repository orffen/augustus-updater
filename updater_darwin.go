// Copyright (c) 2026 Steve Simenic. All rights reserved.
// Use of this source code is governed by a GNU AGPLv3 license
// that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

const (
	dstAppName = "Augustus Unstable.app"
	mountPoint = "/Volumes/AugustusUnstable"
	outFile    = "temp.dmg"
	regex      = `href="([^"]+mac\.dmg)"`
	srcAppName = "augustus.app"
)

var (
	appLocation string
)

func init() {
	home, err := os.UserHomeDir()
	if err != nil {
		fatalError(fmt.Errorf("couldn't find home directory: %w", err))
	}
	appLocation = filepath.Join(home, "Applications", dstAppName)
}

func applyUpdate(url string) error {
	if err := downloadUpdate(url, outFile); err != nil {
		return err
	}
	if err := installFromDMG(outFile); err != nil {
		return err
	}
	if err := os.Remove(outFile); err != nil {
		return fmt.Errorf("couldn't remove temporary file %v: %w", outFile, err)
	}
	return nil
}

func installFromDMG(dmgFile string) error {
	cmd := exec.Command("hdiutil", "attach", dmgFile, "-mountpoint", mountPoint, "-nobrowse", "-quiet")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to mount DMG: %v", err)
	}
	defer func() { _ = exec.Command("hdiutil", "detach", mountPoint, "-quiet").Run() }()
	src := filepath.Join(mountPoint, srcAppName)
	dst := filepath.Dir(appLocation)
	if err := os.MkdirAll(dst, 0755); err != nil {
		return fmt.Errorf("couldn't create applications directory: %w", err)
	}
	_ = os.RemoveAll(appLocation)
	copyCmd := exec.Command("cp", "-R", src, appLocation) // rely on macOS utilities to preserve attributes
	if err := copyCmd.Run(); err != nil {
		return fmt.Errorf("couldn't copy app: %v", err)
	}
	_ = exec.Command("xattr", "-d", "-r", "-s", "com.apple.quarantine", appLocation).Run() // GateKeeper for Augustus Unstable
	return nil
}

func runProgram() {
	cmd := exec.Command("open", appLocation)
	_ = cmd.Start()
	time.Sleep(100 * time.Millisecond)
	os.Exit(0)
}
