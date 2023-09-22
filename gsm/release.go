// Copyright (c) 2023, The GoKi Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gsm

import "fmt"

// Release releases all of the GoKi Go repositories in the current folder,
// recursively updating each one and all of its dependencies, but stopping
// after a couple of iterations due to pseudo-import cycles at the module level.
//
//gti:add
func Release(c *Config) error {
	pkgs, err := ParsePackages()
	if err != nil {
		return fmt.Errorf("error parsing packages: %w", err)
	}
	for _, pkg := range pkgs {
		fmt.Println(pkg.Name)
	}
	return nil
}
