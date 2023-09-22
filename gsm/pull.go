// Copyright (c) 2023, The GoKi Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gsm

import (
	"fmt"
	"io/fs"
	"os"
)

// Pull concurrently pulls all of the Git repositories in the current directory.
//
//gti:add
func Pull(c *Config) error {
	fs.WalkDir(os.DirFS("."), ".", func(path string, d fs.DirEntry, err error) error {
		fmt.Println(path)
		return nil
	})
	return nil
}
