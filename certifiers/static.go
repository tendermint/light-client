package certifiers

import (
	"bytes"

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
	vset    *types.ValidatorSet
	vhash   []byte
}

func NewStatic(chainID string, vals []*types.Validator) StaticCertifier {
	vset := types.NewValidatorSet(vals)
	return StaticCertifier{
		chainID: chainID,
		vset:    vset,
		vhash:   vset.Hash(),
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

	// make sure it has the same validator set we have (static means static)
	if !bytes.Equal(c.vhash, check.Header.ValidatorsHash) {
		return errors.Errorf("Validator hash has changes %X -> %x", c.vhash,
			check.Header.ValidatorsHash)
	}

	// then make sure we have the proper signatures for this
	err = c.vset.VerifyCommit(c.chainID, check.Commit.BlockID,
		check.Header.Height, check.Commit)
	return errors.WithStack(err)
}
