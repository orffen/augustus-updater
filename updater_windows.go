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
	"sync"
	"syscall"
	"time"
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

	var downloadWG sync.WaitGroup
	var unzipWG sync.WaitGroup
	errChan := make(chan error, len(urls))

	for i, u := range urls {
		downloadWG.Add(1)
		unzipWG.Add(1)
		file := fmt.Sprintf(outFile, i)
		go func(downloadURL, file string) {
			defer unzipWG.Done()
			defer func() { _ = os.Remove(file) }()
			if err := downloadUpdate(downloadURL, file); err != nil {
				errChan <- err
				downloadWG.Done()
				return
			}
			downloadWG.Done()
			downloadWG.Wait()
			if err := unzip(context.TODO(), file); err != nil {
				errChan <- err
				return
			}
		}(u, file)
	}
	unzipWG.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			return err
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
