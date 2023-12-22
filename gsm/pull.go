// Copyright (c) 2023, The Goki Authors. All rights reserved.
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

	"goki.dev/xe"
)

// Pull concurrently pulls all of the Git repositories in the current directory.
func Pull(c *Config) error { //gti:add
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
			err := xe.Major().SetDir(dir).Run("git", "pull")
			if err != nil {
				errs = append(errs, fmt.Errorf("error pulling %q: %w", dir, err))
			}
		}()
		return nil
	})
	wg.Wait()
	return errors.Join(errs...)
}
