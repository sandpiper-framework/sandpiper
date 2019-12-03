// Copyright Auto Care Association. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

// package version is used to define product identification information.
// The exported "Version" variable is set by the build process.
package version

import (
	"fmt"
	"strings"
)

// program version is inserted by `go build` task from latest github tag
var Version = "unknown"

const copyright = "Copyright (c) Auto Care Association. All rights reserved."

func Banner() string {
	return fmt.Sprintf("%s\nSandpiper API Server (%s)\n%s\n", product(), Version, copyright)
}

func product() string {
	s := []string{
		"  _____                 _       _",
		" / ____|               | |     (_)",
		"| (___   __ _ _ __   __| |_ __  _ _ __   ___ _ __",
		" \\___ \\ / _` | '_ \\ / _` | '_ \\| | '_ \\ / _ \\ '__|",
		" ____) | (_| | | | | (_| | |_) | | |_) |  __/ |",
		"|_____/ \\__,_|_| |_|\\__,_| .__/|_| .__/ \\___|_|",
		"                         | |     | |",
		"                         |_|     |_|",
	}
	return strings.Join(s,"\n")
}
