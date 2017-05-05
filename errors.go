package lightclient

import (
	"fmt"

	"github.com/pkg/errors"
)

type errHeightMismatch struct {
	h1, h2 int
}

func (e errHeightMismatch) Error() string {
	return fmt.Sprintf("Blocks don't match - %d vs %d", e.h1, e.h2)
}

// IsHeightMismatchErr checks whether an error is due to data from different blocks
func IsHeightMismatchErr(err error) bool {
	if err == nil {
		return false
	}
	_, ok := errors.Cause(err).(errHeightMismatch)
	return ok
}

func ErrHeightMismatch(h1, h2 int) error {
	err := errHeightMismatch{h1, h2}
	return errors.WithStack(err)
}
