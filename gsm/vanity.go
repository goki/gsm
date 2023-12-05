// Copyright (c) 2023, The GoKi Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gsm

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/iancoleman/strcase"
	"goki.dev/xe"
)

type newVanityTmplData struct {
	Name     string
	RepoName string
	Title    string
}

var newVanityTmpl = template.Must(template.New("newVanity").Parse(
	`+++
title = '{{.Title}}'
repo = 'https://github.com/goki/{{.RepoName}}'
packages = ['goki.dev/{{.Name}}']
+++
`))

// NewVanity makes a new vanity import URL page for the config
// repository name. It should only be called in the root directory
// of the goki.github.io repository. It commits and pushes the page.
func NewVanity(c *Config) error { //gti:add
	b := bytes.Buffer{}
	d := newVanityTmplData{Name: c.Repository, RepoName: c.Repository, Title: strcase.ToCamel(c.Repository)}
	// we cut any later parts of the repository name (major version suffixes,
	// submodules, etc), but leave them in the module name
	if before, _, has := strings.Cut(c.Repository, "/"); has {
		d.RepoName = before
		d.Title = strcase.ToCamel(d.RepoName)
	}
	err := newVanityTmpl.Execute(&b, d)
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
	err = xe.Run("git", "add", fname)
	if err != nil {
		return fmt.Errorf("error adding to git: %w", err)
	}
	err = xe.Run("git", "commit", "-am", "added "+c.Repository)
	if err != nil {
		return err
	}
	return xe.Run("git", "push")
}
