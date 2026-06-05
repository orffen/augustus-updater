// Copyright (c) 2026 Steve Simenic. All rights reserved.
// Use of this source code is governed by a GNU AGPLv3 license
// that can be found in the LICENSE file.

//go:generate go run genassets_windows.go
//go:generate goversioninfo

package main

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"io/fs"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

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

func runProgram() error {
	cmd := exec.Command(exe)
	cmd.Dir = "."
	// CREATE_NEW_PROCESS_GROUP = 0x00000200
	// DETACHED_PROCESS = 0x00000008
	cmd.SysProcAttr = &syscall.SysProcAttr{
		CreationFlags: 0x00000200 | 0x00000008,
	}
	err := cmd.Start()
	if err != nil {
		return err
	}
	return nil
}

func unzip(ctx context.Context, file string) error {
	r, err := zip.OpenReader(file)
	if err != nil {
		return fmt.Errorf("couldn't open zip %v: %w", file, err)
	}
	defer func() { _ = r.Close() }()
	return fs.WalkDir(r, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil || path == "." {
			return err
		}
		if d.IsDir() { // handle empty directories
			return os.MkdirAll(path, 0755)
		}
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return err
		}
		if err := ctx.Err(); err != nil {
			return fmt.Errorf("unzip aborted: %w", err)
		}
		src, err := r.Open(path)
		if err != nil {
			return err
		}
		dst, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			_ = src.Close()
			return err
		}
		_, err = io.Copy(dst, src)
		// cleanup first
		_ = dst.Close()
		_ = src.Close()
		if err != nil { // io.Copy error
			return err
		}
		return nil
	})
}
