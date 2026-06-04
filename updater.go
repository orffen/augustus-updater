// Copyright (c) 2026 Steve Simenic. All rights reserved.
// Use of this source code is governed by a GNU AGPLv3 license
// that can be found in the LICENSE file.

package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
)

const (
	updateURL = "https://augustus.josecadete.net/"
)

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

func downloadUpdate(ctx context.Context, url string, filename string) error {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request for %s: %w", url, err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("download failed for %s: %w", url, err)
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
