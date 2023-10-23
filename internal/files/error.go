package files

import (
	"errors"
)

// ErrEntryNotDir is returned if an entry is supposed to be a directory but is
// not.
var ErrEntryNotDir = errors.New("entry is not a directory")

// ErrEntryNotExists is returned if an entry that is looked up does not exist.
var ErrEntryNotExists = errors.New("entry does not exist")

// ErrEntryExists is returned if an entry exists that was not expected.
var ErrEntryExists = errors.New("entry exists")
