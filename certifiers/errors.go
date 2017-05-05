package certifiers

import (
	rawerr "errors"

	"github.com/pkg/errors"
)

var (
	errValidatorsChanged  = rawerr.New("Validators differ between header and certifier")
	errIsSeedNotFoundErr  = rawerr.New("Seed not found by provider")
	errIsTooMuchChangeErr = rawerr.New("Validators change too much to safely update")
	errIsPastTimeErr      = rawerr.New("Update older than certifier height")
	errIsNoPathFoundErr   = rawerr.New("Cannot find a path of validators")
)

// IsSeedNotFoundErr checks whether an error is due to missing data
func IsSeedNotFoundErr(err error) bool {
	return err != nil && (errors.Cause(err) == errIsSeedNotFoundErr)
}

func ErrIsSeedNotFoundErr() error {
	return errors.WithStack(errIsSeedNotFoundErr)
}

// IsValidatorsChangedErr checks whether and error is due
// to a differing validator set
func IsValidatorsChangedErr(err error) bool {
	return err != nil && (errors.Cause(err) == errValidatorsChanged)
}

func ErrIsValidatorsChangedErr() error {
	return errors.WithStack(errValidatorsChanged)
}

// IsTooMuchChangeErr checks whether and error is due to too much change
// between these validators sets
func IsTooMuchChangeErr(err error) bool {
	return err != nil && (errors.Cause(err) == errIsTooMuchChangeErr)
}

func ErrIsTooMuchChangeErr() error {
	return errors.WithStack(errIsTooMuchChangeErr)
}

func IsPastTimeErr(err error) bool {
	return err != nil && (errors.Cause(err) == errIsPastTimeErr)
}

func ErrIsPastTimeErr() error {
	return errors.WithStack(errIsPastTimeErr)
}

// IsNoPathFoundErr checks whether an error is due to no path of
// validators in provider from where we are to where we want to be
func IsNoPathFoundErr(err error) bool {
	return err != nil && (errors.Cause(err) == errIsNoPathFoundErr)
}

func ErrIsNoPathFoundErr() error {
	return errors.WithStack(errIsNoPathFoundErr)
}
