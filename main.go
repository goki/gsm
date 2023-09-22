// Copyright (c) 2023, The GoKi Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

//go:generate goki generate ./...

import (
	"goki.dev/grease"
	"goki.dev/gsm/gsm"
)

func main() {
	opts := grease.DefaultOptions("gsm", "GSM", "CLI and GUI tools for maintaining the source code of GoKi itself (GoKi Source Management)")
	grease.Run(opts, &gsm.Config{}, gsm.NewVanity, gsm.Clone)
}
