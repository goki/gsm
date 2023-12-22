// Copyright (c) 2023, The Goki Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gsm

import (
	"errors"
	"fmt"
	"os"
	"sync"

	"goki.dev/xe"
)

// Clone concurrently clones all of the Goki Go repositories into the current directory.
// It does not clone repositories that the user already has in the current directory.
func Clone(c *Config) error { //gti:add
	reps, err := GetWebsiteRepositories()
	if err != nil {
		return fmt.Errorf("error getting repositories: %w", err)
	}
	wg := sync.WaitGroup{}
	wg.Add(len(reps))
	var errs []error
	for _, rep := range reps {
		rep := rep
		go func() {
			defer wg.Done()
			fi, err := os.Stat(rep.Name)
			if err == nil { // no error means it already exists
				if fi.IsDir() { // if we already have dir, we don't need to clone
					return
				} else {
					errs = append(errs, fmt.Errorf("file %q (for repository %q) already exists and is not a directory", rep.Name, rep.Title))
					return
				}
			}
			err = xe.Run("git", "clone", rep.RepositoryURL)
			if err != nil {
				errs = append(errs, fmt.Errorf("error cloning repository: %w", err))
			}
		}()
	}
	wg.Wait()
	return errors.Join(errs...)
}
