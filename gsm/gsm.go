// Copyright (c) 2023, The GoKi Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package gsm provides functions for maintaining the source code of GoKi itself (GoKi Source Management)
package gsm

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/iancoleman/strcase"
	"goki.dev/xe"
)

// NewVanity makes a new vanity import URL page for the config
// repository name. It should only be called in the root directory
// of the goki.github.io repository.
//
//gti:add
func NewVanity(c *Config) error {
	data := fmt.Sprintf(`---
title: %s
repo: "https://github.com/goki/%s"
packages: ["goki.dev/%s"]
---
`, strcase.ToCamel(c.Repository), c.Repository, c.Repository)
	dir := filepath.Join("content", "en", c.Repository)
	err := os.MkdirAll(dir, 0770)
	if err != nil {
		return fmt.Errorf("error making repository directory: %w", err)
	}
	fname := filepath.Join(dir, "_index.md")
	err = os.WriteFile(fname, []byte(data), 0666)
	if err != nil {
		return fmt.Errorf("error writing to _index.md file for repository: %w", err)
	}
	xc := xe.DefaultConfig()
	xc.Fatal = false
	err = xe.Run(xc, "git", "add", fname)
	if err != nil {
		return fmt.Errorf("error adding to git: %w", err)
	}
	return nil
}
