// Copyright The Sandpiper Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

// Package version is used to define product identification information.
// The exported "Version" variable is set by the build process.
package version

import (
	"fmt"
)

// Version identification is updated by `go build` task from latest github tag
var Version = "unknown"

// Banner prints identifying information about the server.
func Banner() string {
	const copyright = "Copyright 2020 The Sandpiper Authors. All rights reserved."

	return fmt.Sprintf("Sandpiper sync (%s)\n%s\n", Version, copyright)
}