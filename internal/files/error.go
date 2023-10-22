package files

import (
	"errors"
)

var ErrFSEntryNotDir = errors.New("entry is not a directory")
var ErrFSEntryNotExists = errors.New("entry does not exist")
var ErrFSEntryExists = errors.New("entry does not exist")
