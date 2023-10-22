package files

import (
	"errors"
)

var errFSEntryNotDir = errors.New("entry is not a directory")
var errFSEntryNotExists = errors.New("entry does not exist")
var errFSEntryExists = errors.New("entry does not exist")
