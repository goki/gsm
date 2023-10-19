// Copyright (c) 2023, The GoKi Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gsm

import (
	"path/filepath"

	"goki.dev/xe"
)

// InstallTools installs all of the GoKi tools required for development on
// the GoKi codebase (goki, gsm, gtigen, and enumgen). It should be run in a
// directory containing all of the goki repositories (set up with gsm clone),
// and with a go.work file contianing all of those repositories (set up with gsm work).
//
//gti:add
func InstallTools(c *Config) error {
	paths := []string{
		"goki",
		"gsm",
		filepath.Join("gti", "cmd", "gtigen"),
		filepath.Join("enums", "cmd", "enumgen"),
	}
	for _, path := range paths {
		err := xe.Major().SetDir(path).Run("go", "install")
		if err != nil {
			return err
		}
	}
	return nil
}
