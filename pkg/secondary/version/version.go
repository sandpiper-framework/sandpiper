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
	const copyright = "Copyright (c) Auto Care Association. All rights reserved."

	return fmt.Sprintf("%s\nSandpiper Client (%s)\n%s\n", product(), Version, copyright)
}

func product() string {
	// http://patorjk.com/software/taag/#p=display&f=Standard&t=Sandpiper
	// it includes back ticks, which makes this more difficult (replace with `+"`"+`).

	const s = `
  ____                  _       _                 
 / ___|  __ _ _ __   __| |_ __ (_)_ __   ___ _ __ 
 \___ \ / _` + "`" + ` | '_ \ / _` + "`" + ` | '_ \| | '_ \ / _ \ '__|
  ___) | (_| | | | | (_| | |_) | | |_) |  __/ |   
 |____/ \__,_|_| |_|\__,_| .__/|_| .__/ \___|_|   
                         |_|     |_|              
`
	return s
}
