// Copyright (c) 2026 Steve Simenic. All rights reserved.
// Use of this source code is governed by a GNU AGPLv3 license
// that can be found in the LICENSE file.

//go:build ignore

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

type VersionInfo struct {
	FixedFileInfo  any            `json:"FixedFileInfo"`
	StringFileInfo StringFileInfo `json:"StringFileInfo"`
	IconPath       string         `json:"IconPath"`
	ManifestPath   string         `json:"ManifestPath"`
}

type StringFileInfo struct {
	FileDescription  string `json:"FileDescription"`
	OriginalFilename string `json:"OriginalFilename"`
	ProductName      string `json:"ProductName"`
	LegalCopyright   string `json:"LegalCopyright"`
	FileVersion      string `json:"FileVersion"`
	ProductVersion   string `json:"ProductVersion"`
}

func main() {
	ver := "0.0.0-dev" // fallback version string
	major, minor, patch, build := 0, 0, 0, 0

	cmd := exec.Command("git", "describe", "--always")
	if out, err := cmd.Output(); err == nil {
		ver = strings.TrimSpace(string(out))
		re := regexp.MustCompile(`v(\d+)\.(\d+)\.(\d+)`)
		if matches := re.FindStringSubmatch(ver); len(matches) == 4 {
			fmt.Sscanf(matches[1]+" "+matches[2]+" "+matches[3], "%d %d %d", &major, &minor, &patch)
		}
	}

	verInfo := VersionInfo{
		FixedFileInfo: map[string]any{
			"FileVersion":    map[string]int{"Major": major, "Minor": minor, "Patch": patch, "Build": build},
			"ProductVersion": map[string]int{"Major": major, "Minor": minor, "Patch": patch, "Build": build},
		},
		StringFileInfo: StringFileInfo{
			FileDescription:  "Augustus Unstable version updater",
			OriginalFilename: "augustus-updater.exe",
			ProductName:      "Augustus Updater",
			LegalCopyright:   "Copyright (c) 2026 Steve Simenic. Licensed under AGPLv3. App icon (c) Augustus developers under AGPLv3.",
			FileVersion:      strings.TrimPrefix(ver, "v"),
			ProductVersion:   strings.TrimPrefix(ver, "v"),
		},
		IconPath:     "assets/augustus-updater.ico",
		ManifestPath: "assets/augustus-updater.manifest",
	}

	jsonData, _ := json.Marshal(verInfo)
	_ = os.WriteFile("versioninfo.json", jsonData, 0644)
	fmt.Printf("Wrote versioninfo.json for %s", ver)
}
