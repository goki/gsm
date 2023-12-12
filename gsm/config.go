// Copyright (c) 2023, The GoKi Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package gsm provides functions for maintaining the source code of GoKi itself (GoKi Source Management)
package gsm

// Config contains the configuration information for the GSM tool
type Config struct { //gti:add

	// Update is whether to update dependencies and tidy modules
	// when doing a release cycle. It should only be turned off
	// in rare cases in which updating dependencies or tidying
	// modules would cause problems or is not possible.
	Update bool `cmd:"release" def:"true"`

	// The name of the repository to create a vanity import site for.
	// A major version suffix can be added to the end of the repository name
	// (eg: "gi/v2")
	Repository string `cmd:"new-vanity" posarg:"0"`

	// the config info for the make-ios-framework command
	IOSFramework IOSFramework `cmd:"make-ios-framework"`
}

type IOSFramework struct { //gti:add

	// the path of the .dylib file
	Dylib string

	// the name of the resulting framework
	Framework string

	// the name/email address of the developer to have sign the framework
	Developer string

	// the organization to use in the bundle id for the resulting framework
	Organization string
}
