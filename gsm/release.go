// Copyright (c) 2023, The GoKi Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gsm

import (
	"fmt"
	"slices"
	"strings"

	"goki.dev/grog"
	"goki.dev/xe"
)

// Release releases all of the GoKi Go repositories in the current folder with goki.dev
// vanity import URLs (those without vanity import URLs should be released separately),
// recursively updating each one and all of its dependencies, but stopping
// after a couple of iterations due to pseudo-import cycles at the module level.
func Release(c *Config) error { //gti:add
	reps, err := GetLocalRepositories()
	if err != nil {
		return fmt.Errorf("error parsing packages: %w", err)
	}
	repsm := map[string]*Repository{} // map of repositories
	for _, rep := range reps {
		if skipRepo(rep) {
			continue
		}
		repsm[rep.VanityURL] = rep

		tag, err := xe.Minor().SetDir(rep.Name).Output("git", "describe", "--abbrev=0")
		if err != nil {
			return fmt.Errorf("error getting latest tag for repository %q: %w", rep.Name, err)
		}
		rep.Version = tag
		rep.Changed, err = RepositoryHasChanged(rep, tag)
		if err != nil {
			return err
		}

		// if we have GoKi imports, we need to update them first, so we skip
		if len(rep.GoKiImports) > 0 {
			continue
		}

		// don't use sum db to avoid problems (see https://github.com/golang/go/issues/42809)
		xc := xe.Major().SetDir(rep.Name).SetEnv("GONOSUMDB", "*")

		err = xc.Run("go", "get", "-u", "./...")
		if err != nil {
			return fmt.Errorf("error updating deps for repository %q: %w", rep.Name, err)
		}
		err = xc.Run("go", "mod", "tidy")
		if err != nil {
			return fmt.Errorf("error tidying mod for repository %q: %w", rep.Name, err)
		}

		// check again if we are changed after updating deps and mod
		rep.Changed, err = RepositoryHasChanged(rep, tag)
		if err != nil {
			return err
		}

		if rep.Changed { // if we are changed and have no GoKi imports, we can release right now
			err := ReleaseRepository(rep)
			if err != nil {
				return err
			}
			rep.Released = true
		}
	}

	// now that we have done a first pass to get the ones with no GoKi imports,
	// we check again based on the newly released ones, until we run out of
	// repositories to release. we set a backup break point of 10.
	for i := 0; i < 10; i++ {
		needRelease := false // whether we still have something that needs to be released but can't be
		for _, rep := range reps {
			if skipRepo(rep) {
				continue
			}
			if rep.Released { // if we are already released, we skip
				continue
			}
			hasGoKiImport := false // whether we still have changed but unreleased GoKi imports

			// don't use sum db to avoid problems (see https://github.com/golang/go/issues/42809)
			xc := xe.Major().SetDir(rep.Name).SetEnv("GONOSUMDB", "*")

			for _, imp := range rep.GoKiImports {
				impr := repsm[imp]
				if impr == nil {
					return fmt.Errorf("missing repository for import %q; you might need to run gsm clone", imp)
				}
				if !impr.Changed { // if the import hasn't been changed, we don't need to update it
					continue
				}
				if !impr.Released { // if the import has changed but hasn't been released, we have to wait for them to release first
					hasGoKiImport = true
					continue
				}
				// otherwise, we need to update to the latest release
				err := xc.Run("go", "get", impr.VanityURL+"@"+impr.Version)
				if err != nil {
					return fmt.Errorf("error updating GoKi import %q for repository %q: %w", impr.Name, rep.Name, err)
				}
			}
			// we skip if we still have unreleased GoKi imports,
			// unless we are on the second pass and are one of the three
			// special cyclically importing repositories
			if hasGoKiImport && !(i == 1 && (rep.Name == "enums" || rep.Name == "gti" || rep.Name == "grease")) {
				needRelease = true
				continue
			}

			// now we make sure we have the latest versions of everything
			err := xc.Run("go", "get", "-u", "./...")
			if err != nil {
				return fmt.Errorf("error updating deps for repository %q: %w", rep.Name, err)
			}
			err = xc.Run("go", "mod", "tidy")
			if err != nil {
				return fmt.Errorf("error tidying mod for repository %q: %w", rep.Name, err)
			}
			tag, err := xe.Minor().SetDir(rep.Name).Output("git", "describe", "--abbrev=0")
			if err != nil {
				return fmt.Errorf("error getting latest tag for repository %q: %w", rep.Name, err)
			}
			rep.Version = tag

			// we skip if we still haven't changed
			rep.Changed, err = RepositoryHasChanged(rep, rep.Version)
			if err != nil {
				return err
			}
			if !rep.Changed {
				continue
			}

			// otherwise, we can release
			err = ReleaseRepository(rep)
			if err != nil {
				return err
			}
			rep.Released = true
		}
		if !needRelease {
			break
		}
	}
	return nil
}

// skipRepo returns whether to skip the given repository.
// TODO(kai): remove this TEMPORARY fix for some repos being a WIP
func skipRepo(rep *Repository) bool {
	skips := []string{"gide", "gipy", "grid", "gopix", "gosl", "gear", "goki.github.io", "rqlite", "gorqlite"}
	return slices.Contains(skips, rep.Name)
}

// RepositoryHasChanged returns whether the given repository
// has changed since the given Git version tag.
func RepositoryHasChanged(rep *Repository, tag string) (bool, error) {
	diff, err := xe.Minor().SetDir(rep.Name).Output("git", "diff", tag)
	if err != nil {
		return false, fmt.Errorf("error getting diff from latest tag %q for repository %q: %w", tag, rep.Name, err)
	}
	return diff != "", nil
}

// ReleaseRepository releases the given repository by calling
// "goki update-version" and "goki release".
func ReleaseRepository(rep *Repository) error {
	xc := xe.Major().SetDir(rep.Name)

	err := xc.Run("goki", "update-version")
	if err != nil {
		return fmt.Errorf("error getting updating version of repository %q: %w", rep.Name, err)
	}
	nv, err := xc.Output("goki", "get-version")
	if err != nil {
		return fmt.Errorf("error getting new version of repository %q: %w", rep.Name, err)
	}
	// we only want the part before the newline (the version)
	rep.Version, _, _ = strings.Cut(nv, "\n")
	err = xc.Run("goki", "release")
	if err != nil {
		return fmt.Errorf("error releasing repository %q: %w", rep.Name, err)
	}
	grog.PrintlnWarn(grog.SuccessColor("Released "), grog.CmdColor(rep.Name))
	return nil
}
