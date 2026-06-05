// Copyright (c) 2026 Steve Simenic. All rights reserved.
// Use of this source code is governed by a GNU AGPLv3 license
// that can be found in the LICENSE file.

//go:build ignore

package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"text/template"
)

type InfoPlistData struct {
	BundleExecutable       string
	BundleIdentifier       string
	BundleVersion          string
	HumanReadableCopyright string
	IconFile               string
}

const InfoPlistTemplate = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>CFBundleExecutable</key>
	<string>{{.BundleExecutable}}</string>
	<key>CFBundleIdentifier</key>
	<string>{{.BundleIdentifier}}</string>
	<key>CFBundlePackageType</key>
	<string>APPL</string>
	<key>CFBundleShortVersionString</key>
	<string>{{.BundleVersion}}</string>
	<key>NSHumanReadableCopyright</key>
	<string>{{.HumanReadableCopyright}}</string>
	<key>CFBundleIconFile</key>
	<string>{{.IconFile}}</string>
</dict>
</plist>
`

func main() {
	ver := "0.0.0-dev" // fallback version string

	if out, err := exec.Command("git", "describe", "--tags", "--abbrev=0").Output(); err == nil {
		ver = strings.TrimSpace(string(out))
		ver = strings.TrimPrefix(ver, "v")
	}

	data := InfoPlistData{
		BundleExecutable:       "augustus-updater-mac",
		BundleIdentifier:       "com.github.orffen.augustus-updater",
		BundleVersion:          ver,
		HumanReadableCopyright: "Copyright (c) 2026 Steve Simenic. Licensed under AGPLv3. App icon (c) Augustus developers under AGPLv3.",
		IconFile:               "augustus-updater.icns",
	}

	tmpl, err := template.New("infoPlist").Parse(InfoPlistTemplate)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing template: %v\n", err)
		os.Exit(1)
	}
	file, err := os.Create("Info.plist")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating Info.plist: %v\n", err)
		os.Exit(1)
	}
	defer func() { _ = file.Close() }()
	err = tmpl.Execute(file, data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing template: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Wrote Info.plist for %s\n", ver)
}
