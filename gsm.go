// Copyright (c) 2023, The Goki Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"goki.dev/grease"
	"goki.dev/gsm/cmd"
)

func main() {
	opts := grease.DefaultOptions("gsm", "GSM", "CLI and GUI tools for maintaining the source code of Goki itself (Goki Source Management)")
	grease.Run(opts, &cmd.Config{}, cmd.Clone, cmd.Pull, cmd.Changed, cmd.Release, cmd.Work, cmd.InstallTools, cmd.Gendex, cmd.NewVanity, cmd.MakeIOSFramework)
}
