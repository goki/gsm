// Copyright (c) 2023, The GoKi Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package gsm provides functions for maintaining the source code of GoKi itself (GoKi Source Management)
package gsm

import "fmt"

// NewVanity makes a new vanity import URL page for the config
// repository name. It should only be called in the root directory
// of the goki.github.io repository.
//
//gti:add
func NewVanity(c *Config) error {
	fmt.Println("making new vanity for", c.Package)
	return nil
}
