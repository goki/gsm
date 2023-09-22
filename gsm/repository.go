// Copyright (c) 2023, The GoKi Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gsm

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path"
	"sync"

	"goki.dev/xe"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// Clone concurrently clones all of the GoKi Go repositories into the current directory.
// It does not clone repositories that the user already has in the current directory.
//
//gti:add
func Clone(c *Config) error {
	reps, err := GetRepositories()
	if err != nil {
		return fmt.Errorf("error getting repositories: %w", err)
	}
	wg := sync.WaitGroup{}
	wg.Add(len(reps))
	var errs []error
	for _, rep := range reps {
		rep := rep
		go func() {
			defer wg.Done()
			fi, err := os.Stat(rep.Name)
			if err == nil { // no error means it already exists
				if fi.IsDir() { // if we already have dir, we don't need to clone
					return
				} else {
					errs = append(errs, fmt.Errorf("file %q (for repository %q) already exists and is not a directory", rep.Name, rep.Title))
					return
				}
			}
			xc := xe.DefaultConfig()
			xc.Fatal = false
			err = xe.Run(xc, "git", "clone", rep.RepositoryURL)
			if err != nil {
				errs = append(errs, fmt.Errorf("error cloning repository: %w", err))
				return
			}
		}()
	}
	wg.Wait()
	return errors.Join(errs...)
}

// Repository represents a GoKi Go repository
type Repository struct {
	// The actual GitHub name of the repository
	Name string
	// The formatted title of the repository
	Title string
	// The URL of the GitHub repository (including https://)
	RepositoryURL string
	// The goki.dev vanity import URL of the repository (not including https://)
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
	reps, err := extractRepositories(tree)
	if err != nil {
		return nil, fmt.Errorf("error extracting repositories from HTML: %w", err)
	}
	return reps, nil
}

// extractRepositories extracts repositories from the given HTML node
// that should be the root node of the https://goki.dev/repositories page.
func extractRepositories(node *html.Node) ([]*Repository, error) {
	nodes := appendAll(nil, node, func(n *html.Node) bool {
		if n.DataAtom != atom.Div {
			return false
		}
		for _, a := range n.Attr {
			if a.Key == "class" && a.Val == "entry" {
				return true
			}
		}
		return false
	})
	res := []*Repository{}
	for _, n := range nodes {
		// structure:
		// <div class=entry>
		// 		<h5>
		// 			<a href="/name/">Title></a>
		//		</h5>
		// </div>
		if n.FirstChild == nil || n.FirstChild.FirstChild == nil {
			return nil, fmt.Errorf("got nil first child or first child's child for entry div node %#v", n)
		}
		a := n.FirstChild.FirstChild
		href := ""
		for _, attr := range a.Attr {
			if attr.Key == "href" {
				href = attr.Val
			}
		}
		if href == "" {
			return nil, fmt.Errorf("could not get href for node %#v", a)
		}
		if a.FirstChild.Type != html.TextNode {
			return nil, fmt.Errorf("expected text node as first child of node %#v", a)
		}
		rep := &Repository{
			Name:  path.Base(href),
			Title: a.FirstChild.Data,
		}
		rep.RepositoryURL = "https://github.com/goki/" + rep.Name
		rep.VanityURL = "goki.dev/" + rep.Name
		res = append(res, rep)
	}
	return res, nil
}

// matchFunc matches HTML nodes.
type matchFunc func(*html.Node) bool

// appendAll recursively traverses the parse tree rooted under the provided
// node and appends all nodes matched by the matchFunc to dst.
func appendAll(dst []*html.Node, n *html.Node, mf matchFunc) []*html.Node {
	if mf(n) {
		dst = append(dst, n)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		dst = appendAll(dst, c, mf)
	}
	return dst
}
