// Copyright (c) 2023, The GoKi Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gsm

import (
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"sync"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

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
	// The imports of the repository
	Imports []string
}

// GetLocalRepositories concurrently gets all of the GoKi Go
// repositories in the current directory on the local filesystem.
func GetLocalRepositories() ([]*Repository, error) {
	wg := sync.WaitGroup{}
	errs := []error{}
	res := []*Repository{}
	fs.WalkDir(os.DirFS("."), ".", func(path string, d fs.DirEntry, err error) error {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if d.Name() != "go.mod" {
				return
			}
			dir := filepath.Dir(path)
			mod, err := os.ReadFile(path)
			if err != nil {
				errs = append(errs, fmt.Errorf("error reading mod file for %q: %w", dir, err))
			}
			fmt.Println(string(mod))
		}()
		return nil
	})
	wg.Wait()
	return res, errors.Join(errs...)
}

// GetWebsiteRepositories gets all of the GoKi Go repositories as [Repository]
// objects from the https://goki.dev/repositories page.
func GetWebsiteRepositories() ([]*Repository, error) {
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
