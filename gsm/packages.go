// Copyright (c) 2023, The GoKi Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gsm

import (
	"golang.org/x/tools/go/packages"
)

// ParsePackages parses all of the packages in the current
// directory (assumed to be the GoKi root directory).
func ParsePackages() ([]*packages.Package, error) {
	pcfg := &packages.Config{
		Mode:  PackageModes(),
		Tests: false,
	}
	pkgs, err := packages.Load(pcfg, "./...")
	if err != nil {
		return nil, err
	}
	return pkgs, err
}

// PackageModes returns the package load modes needed for gsm.
func PackageModes() packages.LoadMode {
	// TODO: do we need packages.NeedDeps?
	res := packages.NeedName | packages.NeedImports
	return res
}
