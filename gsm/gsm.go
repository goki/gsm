// Copyright (c) 2023, The GoKi Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package gsm provides functions for maintaining the source code of GoKi itself (GoKi Source Management)
package gsm

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/iancoleman/strcase"
	"goki.dev/xe"
)

type newVanityTmplData struct {
	Name  string
	Title string
}

var newVanityTmpl = template.Must(template.New("newVanity").Parse(
	`---
title: {{.Title}}
repo: "https://github.com/goki/{{.Name}}"
packages: ["goki.dev/{{.Name}}"]
---
`))

// NewVanity makes a new vanity import URL page for the config
// repository name. It should only be called in the root directory
// of the goki.github.io repository.
//
//gti:add
func NewVanity(c *Config) error {
	b := bytes.Buffer{}
	err := newVanityTmpl.Execute(&b, newVanityTmplData{Name: c.Repository, Title: strcase.ToCamel(c.Repository)})
	if err != nil {
		return fmt.Errorf("programmer error: error executing vanity URL file template: %w", err)
	}
	dir := filepath.Join("content", "en", c.Repository)
	err = os.MkdirAll(dir, 0770)
	if err != nil {
		return fmt.Errorf("error making vanity URL directory: %w", err)
	}
	fname := filepath.Join(dir, "_index.md")
	err = os.WriteFile(fname, b.Bytes(), 0666)
	if err != nil {
		return fmt.Errorf("error writing to _index.md file for vanity URL: %w", err)
	}
	xc := xe.DefaultConfig()
	xc.Fatal = false
	err = xe.Run(xc, "git", "add", fname)
	if err != nil {
		return fmt.Errorf("error adding to git: %w", err)
	}
	return nil
}
