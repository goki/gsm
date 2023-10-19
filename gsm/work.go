// Copyright (c) 2023, The GoKi Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gsm

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"goki.dev/glop/dirs"
	"goki.dev/xe"
)

// Work adds all of the Go modules in the current directory to the go.work
// file in the current directory.
//
//gti:add
func Work(c *Config) error {
	ex, err := dirs.FileExists("go.work")
	if err != nil {
		return err
	}
	if !ex {
		err := xe.Run("go", "work", "init")
		if err != nil {
			return err
		}
	}
	return fs.WalkDir(os.DirFS("."), ".", func(path string, d fs.DirEntry, err error) error {
		if d.Name() != "go.mod" {
			return nil
		}
		dir := filepath.Dir(path)
		// TODO: figure out a more sustainable solution to this temporary workaround
		if dir == "gipy" || dir == "goki.github.io" || dir == "android-go" || strings.Contains(dir, "internal") {
			return nil
		}
		return xe.Run("go", "work", "use")
	})
}
