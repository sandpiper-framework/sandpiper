package sandpiper

// todo: move this error to each test file and remove this file

import (
	"errors"
)

var (
	// ErrGeneric is used for testing purposes and for errors handled later in the callstack
	ErrGeneric = errors.New("generic error")
)
