// Copyright (c) 2023, The Goki Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import (
	"path/filepath"

	"goki.dev/xe"
)

// Gendex runs goki.dev/goki/mobile/gendex.go and install-tools.
// It should be run in the base goki directory whenever
// goki.dev/goosi/driver/android/GoNativeActivty.java is updated.
func Gendex(c *Config) error { //gti:add
	err := xe.Major().SetDir(filepath.Join("goki", "mobile")).Run("go", "generate")
	if err != nil {
		return err
	}
	return InstallTools(c)
}
