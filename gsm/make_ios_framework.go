// Copyright (c) 2023, The GoKi Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gsm

import (
	"bytes"
	"os"
	"text/template"

	"goki.dev/xe"
)

// MakeIOSFramework makes a .framework file for iOS from
// a .dylib file, using the given config information.
//
//gti:add
func MakeIOSFramework(c *Config) error {
	// based on https://stackoverflow.com/a/57795040
	err := xe.Run("install_name_tool", "-id", "@executable_path/"+c.IOSFramework.Framework+".framework/"+c.IOSFramework.Framework, c.IOSFramework.Dylib)
	if err != nil {
		return err
	}
	err = xe.Run("lipo", "-create", c.IOSFramework.Dylib, "-output", c.IOSFramework.Framework)
	if err != nil {
		return err
	}
	err = xe.MkdirAll(c.IOSFramework.Framework+".framework", 0750)
	if err != nil {
		return err
	}
	err = xe.Run("mv", c.IOSFramework.Framework, c.IOSFramework.Framework+".framework")
	if err != nil {
		return err
	}
	buf := &bytes.Buffer{}
	err = plistTmpl.Execute(buf, c)
	if err != nil {
		return err
	}
	err = os.WriteFile(c.IOSFramework.Framework+".framework/Info.plist", buf.Bytes(), 0666)
	if err != nil {
		return err
	}
	err = xe.Run("codesign", "--force", "--dep", "--verbose=2", "--sign", c.IOSFramework.Developer, c.IOSFramework.Framework+".framework")
	if err != nil {
		return err
	}
	return xe.Run("codesign", "-vvvv", c.IOSFramework.Framework+".framework")
}

var plistTmpl = template.Must(template.New("plist").Parse(
	`<?xml version="1.0" encoding="UTF-8"?>
	<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
	<plist version="1.0">
	<dict>
		<key>BuildMachineOSBuild</key>
		<string>22F82</string>
		<key>CFBundleDevelopmentRegion</key>
		<string>en</string>
		<key>CFBundleExecutable</key>
		<string>{{.IOSFramework.Framework}}</string>
		<key>CFBundleIdentifier</key>
		<string>com.{{.IOSFramework.Organization}}.{{.IOSFramework.Framework}}</string>
		<key>CFBundleInfoDictionaryVersion</key>
		<string>6.0</string>
		<key>CFBundleName</key>
		<string>{{.IOSFramework.Framework}}</string>
		<key>CFBundlePackageType</key>
		<string>FMWK</string>
		<key>CFBundleShortVersionString</key>
		<string>1.0</string>
		<key>CFBundleSupportedPlatforms</key>
		<array>
			<string>iPhoneOS</string>
		</array>
		<key>CFBundleVersion</key>
		<string>1</string>
	</dict>
	</plist>
	`))
