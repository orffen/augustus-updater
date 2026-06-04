// Copyright (c) 2026 Steve Simenic. All rights reserved.
// Use of this source code is governed by a GNU AGPLv3 license
// that can be found in the LICENSE file.

package main

import (
	"archive/zip"
	"bufio"
	"context"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	updateURL   = "https://augustus.josecadete.net/"
	versionFile = "download_url.txt"
)

func downloadUpdate(url string, filename string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("download failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad server response from %s: %s", url, resp.Status)
	}
	out, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("couldn't create temp file %v: %w", filename, err)
	}
	defer func() { _ = out.Close() }()
	if _, err := io.Copy(out, resp.Body); err != nil {
		return fmt.Errorf("failed to save download: %w", err)
	}
	return nil
}

func fatalError(v ...any) {
	args := append([]any{"Fatal Error:"}, v...)
	_ = writeVersion(fmt.Sprintf("%s\nWill redownload when next run.", args)) // force a re-download next run
	fmt.Fprintln(os.Stderr, args...)
	fmt.Println("Press ENTER key to quit...")
	reader := bufio.NewScanner(os.Stdin)
	_ = reader.Scan()
	os.Exit(1)
}

func getDownloadURL(regex string) (string, error) {
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
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", fmt.Errorf("local version file unreadable: %w", err)
	}
	return strings.TrimSpace(string(ver)), nil
}

func showError(v ...any) {
	args := append([]any{"Warning:"}, v...)
	fmt.Fprintln(os.Stderr, args)
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

func writeVersion(ver string) error {
	data := []byte(ver)
	return os.WriteFile(versionFile, data, 0644)
}
