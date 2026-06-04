// Copyright (c) 2026 Steve Simenic. All rights reserved.
// Use of this source code is governed by a GNU AGPLv3 license
// that can be found in the LICENSE file.

//go:generate goversioninfo

package main

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
)

const (
	exe     = "./augustus.exe"
	outFile = "update%d.zip"
	regex   = `href="([^"]+windows-64bit\.zip)"`
)

func applyUpdate(updateURL string) error {
	assetsURL, err := url.Parse(updateURL)
	if err != nil {
		return fmt.Errorf("invalid update url: %w", err)
	}
	replacer := strings.NewReplacer(
		"augustus", "assets",
		"windows", "development",
		"-64bit", "",
	)
	assetsURL.Path = replacer.Replace(assetsURL.Path)
	urls := []string{updateURL, assetsURL.String()}
	files := make([]string, len(urls))
	for i := range urls {
		files[i] = fmt.Sprintf(outFile, i)
	}
	defer func() {
		for _, file := range files {
			_ = os.Remove(file)
		}
	}()

	dlGrp, dlCtx := errgroup.WithContext(context.Background())
	for i, u := range urls {
		dlGrp.Go(func() error {
			return downloadUpdate(dlCtx, u, files[i])
		})
	}
	if err := dlGrp.Wait(); err != nil {
		return err
	}

	unzipGrp, unzipCtx := errgroup.WithContext(context.Background())
	for _, file := range files {
		unzipGrp.Go(func() error {
			return unzip(unzipCtx, file)
		})
	}
	return unzipGrp.Wait()
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
