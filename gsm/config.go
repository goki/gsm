// Copyright (c) 2023, The GoKi Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package gsm provides functions for maintaining the source code of GoKi itself (GoKi Source Management)
package gsm

// Config contains the configuration information for the GSM tool
//
//gti:add
type Config struct {

	// the name of the repository to create a vanity import site for
	Repository string `cmd:"new-vanity" posarg:"0"`

	// the config info for the make-ios-framework command
	IOSFramework IOSFramework `cmd:"make-ios-framework"`
}

//gti:add
type IOSFramework struct {

	// the path of the .dylib file
	Dylib string

	// the name of the resulting framework
	Framework string

	// the name/email address of the developer to have sign the framework
	Developer string

	// the organization to use in the bundle id for the resulting framework
	Organization string
}
