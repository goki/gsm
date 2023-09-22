// Copyright (c) 2023, The GoKi Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gsm

import (
	"fmt"
	"net/http"

	"golang.org/x/net/html"
)

// Repository represents a GoKi Go repository
type Repository struct {
	// The actual GitHub name of the repository
	Name string
	// The formatted title of the repository
	Title string
	// The URL of the GitHub repository
	RepositoryURL string
	// The https://goki.dev vanity import URL of the repository
	VanityURL string
}

// GetRepositories gets all of the GoKi Go repositories as [Repository]
// objects from the https://goki.dev/repositories page.
func GetRepositories() ([]*Repository, error) {
	resp, err := http.Get("https://goki.dev/repositories")
	if err != nil {
		return nil, fmt.Errorf("error getting goki.dev/repositories page: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("got status code %d from goki.dev/repositories (expected 200)", resp.StatusCode)
	}
	tree, err := html.Parse(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error parsing HTML from repositories page: %w", err)
	}
	fmt.Println(tree)
	return nil, nil
}

// Clone clones all of the GoKi Go repositories into the current directory.
// It does not clone repositories that the user already has in the current directory.
//
//gti:add
func Clone(c *Config) error {
	reps, err := GetRepositories()
	if err != nil {
		return fmt.Errorf("error getting repositories: %w", err)
	}
	fmt.Println(reps)
	return nil
}
