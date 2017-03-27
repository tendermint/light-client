package certifiers

import (
	rawerr "errors"

	"github.com/pkg/errors"
)

var (
	errValidatorsChanged = rawerr.New("Validators differ between header and certifier")
	errSeedNotFound      = rawerr.New("Seed not found by provider")
	errTooMuchChange     = rawerr.New("Validators change too much to safely update")
	errPastTime          = rawerr.New("Update older than certifier height")
	errNoPathFound       = rawerr.New("Cannot find a path of validators")
)

// SeedNotFound asserts whether an error is due to missing data
func SeedNotFound(err error) bool {
	return err != nil && (errors.Cause(err) == errSeedNotFound)
}

func ErrSeedNotFound() error {
	return errors.WithStack(errSeedNotFound)
}

// ValidatorsChanged asserts whether and error is due
// to a differing validator set
func ValidatorsChanged(err error) bool {
	return err != nil && (errors.Cause(err) == errValidatorsChanged)
}

func ErrValidatorsChanged() error {
	return errors.WithStack(errValidatorsChanged)
}

// TooMuchChange asserts whether and error is due to too much change
// between these validators sets
func TooMuchChange(err error) bool {
	return err != nil && (errors.Cause(err) == errTooMuchChange)
}

func ErrTooMuchChange() error {
	return errors.WithStack(errTooMuchChange)
}

func PastTime(err error) bool {
	return err != nil && (errors.Cause(err) == errPastTime)
}

func ErrPastTime() error {
	return errors.WithStack(errPastTime)
}

// NoPathFound asserts whether an error is due to no path of
// validators in provider from where we are to where we want to be
func NoPathFound(err error) bool {
	return err != nil && (errors.Cause(err) == errNoPathFound)
}

func ErrNoPathFound() error {
	return errors.WithStack(errNoPathFound)
}
