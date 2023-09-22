// Copyright (c) 2023, The GoKi Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gsm

import (
	"fmt"

	"goki.dev/xe"
)

// Release releases all of the GoKi Go repositories in the current folder with goki.dev
// vanity import URLs (those without vanity import URLs should be released separately),
// recursively updating each one and all of its dependencies, but stopping
// after a couple of iterations due to pseudo-import cycles at the module level.
//
//gti:add
func Release(c *Config) error {
	reps, err := GetLocalRepositories()
	if err != nil {
		return fmt.Errorf("error parsing packages: %w", err)
	}
	for _, rep := range reps {
		vc := xe.VerboseConfig()
		vc.Dir = rep.Name
		err := xe.Run(vc, "go", "get", "-u", "./...")
		if err != nil {
			return fmt.Errorf("error updating deps for repository %q: %w", rep.Name, err)
		}

		ec := xe.ErrorConfig()
		ec.Dir = rep.Name
		tag, err := xe.Output(ec, "git", "describe", "--abbrev=0")
		if err != nil {
			return fmt.Errorf("error getting latest tag for repository %q: %w", rep.Name, err)
		}

		diff, err := xe.Output(ec, "git", "diff", tag)
		if err != nil {
			return fmt.Errorf("error getting diff from latest tag %q for repository %q: %w", tag, rep.Name, err)
		}
		if diff == "" { // unchanged, so no release needed
			continue
		}

		if len(rep.GoKiImports) == 0 { // if we have no GoKi imports, we can release right now
			err := xe.Run(vc, "goki", "set-version", tag)
			if err != nil {
				return fmt.Errorf("error getting setting version of repo %q to %q: %w", rep.Name, tag, err)
			}
		}
	}
	return nil
}
