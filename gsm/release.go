// Copyright (c) 2023, The GoKi Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gsm

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
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
		if rep.Name == "gi" || rep.Name == "gi3d" || rep.Name == "gide" || rep.Name == "gipy" || rep.Name == "grid" || rep.Name == "gopix" || rep.Name == "greasi" { // TODO: remove this TEMPORARY fix for some repos being a WIP
			continue
		}
		vc := xe.VerboseConfig()
		vc.Dir = rep.Name
		err := xe.Run(vc, "go", "get", "-u", "./...")
		if err != nil {
			return fmt.Errorf("error updating deps for repository %q: %w", rep.Name, err)
		}
		err = xe.Run(vc, "go", "mod", "tidy")
		if err != nil {
			return fmt.Errorf("error tidying mod for repository %q: %w", rep.Name, err)
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
		rep.Changed = diff != ""
		if !rep.Changed { // unchanged, so no release needed
			continue
		}

		ver, err := semver.NewVersion(tag)
		if err != nil {
			return fmt.Errorf("error getting semver version of repository %q from tag %q: %w", rep.Name, tag, err)
		}
		rep.Version = ver

		if len(rep.GoKiImports) == 0 { // if we have no GoKi imports, we can release right now
			err := ReleaseRepository(rep)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// ReleaseRepository releases the given repository by incrementing its
// patch version and calling "goki set-version" and "goki release".
func ReleaseRepository(rep *Repository) error {
	vc := xe.VerboseConfig()
	vc.Dir = rep.Name

	*rep.Version = rep.Version.IncPatch()

	nver := "v" + rep.Version.String()
	err := xe.Run(vc, "goki", "set-version", nver)
	if err != nil {
		return fmt.Errorf("error getting setting version of repository %q to %q: %w", rep.Name, nver, err)
	}
	err = xe.Run(vc, "goki", "release")
	if err != nil {
		return fmt.Errorf("error releasing repository %q with version %q: %w", rep.Name, nver, err)
	}
	return nil
}
