package certifiers

import (
	"github.com/pkg/errors"
	lc "github.com/tendermint/light-client"
	"github.com/tendermint/tendermint/types"
)

// StaticCertifier assumes a static set of validators, set on
// initilization and checks against them.
//
// Good for testing or really simple chains.  You will want a
// better implementation when the validator set can actually change.
type StaticCertifier struct {
	chainID string
	vals    *types.ValidatorSet
}

func NewStatic(chainID string, vals []*types.Validator) StaticCertifier {
	return StaticCertifier{
		chainID: chainID,
		vals:    types.NewValidatorSet(vals),
	}
}

func (c StaticCertifier) assertCertifier() lc.Certifier {
	return c
}

func (c StaticCertifier) Certify(check lc.Checkpoint) error {
	// do basic sanity checks
	err := check.ValidateBasic(c.chainID)
	if err != nil {
		return err
	}

	// then make sure we have the proper signatures for this
	err = c.vals.VerifyCommit(c.chainID, check.Commit.BlockID,
		check.Header.Height, check.Commit)
	return errors.WithStack(err)
}
