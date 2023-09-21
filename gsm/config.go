// Copyright (c) 2023, The GoKi Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gsm

// Config contains the configuration information for the GSM tool
type Config struct {

	// the name of the repository to create a vanity import site for
	Repository string `cmd:"new-vanity" posarg:"0" desc:"the name of the repository to create a vanity import site for"`
}
