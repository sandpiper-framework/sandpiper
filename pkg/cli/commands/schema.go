// Copyright The Sandpiper Authors. All rights reserved.
// This file is licensed under the Artistic License 2.0.
// License text can be found in the project's LICENSE file.

// Package command implements `sandpiper` commands (add, pull, list, ...)
package command

import (
	"fmt"

	args "github.com/urfave/cli/v2"

	"github.com/sandpiper-framework/sandpiper/pkg/shared/database"
)

// Schema displays current database schema as defined in the migrate code
func Schema(c *args.Context) error {
	fmt.Println(database.Schema())
	return nil
}
