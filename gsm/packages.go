// Copyright (c) 2023, The GoKi Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gsm

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"golang.org/x/tools/go/packages"
)

// ParsePackages parses all of the root packages in the modules
// in the current directory (assumed to be the GoKi root directory).
func ParsePackages() ([]*packages.Package, error) {
	res := []*packages.Package{}
	pcfg := &packages.Config{
		Mode:  PackageModes(),
		Tests: false,
	}
	// need to get separately because [os.DirFS] on a relative path changes when we call [os.Chdir] later
	cdir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("error getting current working directory: %w", err)
	}
	err = fs.WalkDir(os.DirFS(cdir), ".", func(path string, d fs.DirEntry, perr error) error {
		if d.Name() != "go.mod" {
			return nil
		}
		dir := filepath.Dir(path)
		// need to use [os.Chdir] because you must be in the same module for [packages.Load] to work
		err := os.Chdir(dir)
		if err != nil {
			return fmt.Errorf("error changing to module directory %q: %w", dir, err)
		}
		pkgs, err := packages.Load(pcfg, ".")
		if err != nil {
			return fmt.Errorf("error loading package %q: %w", dir, err)
		}
		res = append(res, pkgs...)
		err = os.Chdir(cdir)
		if err != nil {
			return fmt.Errorf("error changing back to base directory %q: %w", cdir, err)
		}
		return nil
	})
	return res, err
}

// PackageModes returns the package load modes needed for gsm.
func PackageModes() packages.LoadMode {
	// TODO: do we need packages.NeedDeps?
	res := packages.NeedName | packages.NeedImports
	return res
}
