// Copyright (c) 2023, The GoKi Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gsm

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"goki.dev/xe"
)

// Changed prints all of the repositories that have been changed
// and need to be updated in version control.
//
//gti:add
func Changed(c *Config) error {
	return fs.WalkDir(os.DirFS("."), ".", func(path string, d fs.DirEntry, perr error) error {
		if d.Name() != ".git" {
			return nil
		}
		dir := filepath.Dir(path)
		out, err := xe.Output(xe.ErrorConfig(), "git", "-C", dir, "diff")
		if err != nil {
			return fmt.Errorf("error getting diff of %q: %w", dir, err)
		}
		if out != "" { // if we have a diff, we are changed
			fmt.Println(dir)
		}
		return nil
	})
}
