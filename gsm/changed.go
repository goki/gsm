// Copyright (c) 2023, The GoKi Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gsm

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sync"

	"goki.dev/grease"
	"goki.dev/xe"
)

// Changed concurrently prints all of the repositories that have been changed
// and need to be updated in version control.
//
//gti:add
func Changed(c *Config) error {
	wg := sync.WaitGroup{}
	errs := []error{}
	fs.WalkDir(os.DirFS("."), ".", func(path string, d fs.DirEntry, err error) error {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if d.Name() != ".git" {
				return
			}
			dir := filepath.Dir(path)
			out, err := xe.Output(xe.ErrorConfig(), "git", "-C", dir, "diff")
			if err != nil {
				errs = append(errs, fmt.Errorf("error getting diff of %q: %w", dir, err))
			}
			if out != "" { // if we have a diff, we have been changed
				fmt.Println(grease.CmdColor(dir))
			}
		}()
		return nil
	})
	wg.Wait()
	fmt.Println("")
	return errors.Join(errs...)
}
