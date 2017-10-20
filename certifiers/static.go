package certifiers

import (
	"bytes"

	"github.com/pkg/errors"

	"github.com/tendermint/tendermint/types"

	certerr "github.com/tendermint/light-client/certifiers/errors"
)

var _ Certifier = &Static{}

// Static assumes a static set of validators, set on
// initilization and checks against them.
//
// Good for testing or really simple chains.  You will want a
// better implementation when the validator set can actually change.
type Static struct {
	chainID string
	VSet    *types.ValidatorSet
	vhash   []byte
}

func NewStatic(chainID string, vals *types.ValidatorSet) *Static {
	return &Static{
		chainID: chainID,
		VSet:    vals,
	}
}

func (c *Static) ChainID() string {
	return c.chainID
}

func (c *Static) Hash() []byte {
	if len(c.vhash) == 0 {
		c.vhash = c.VSet.Hash()
	}
	return c.vhash
}

func (c *Static) Certify(check Checkpoint) error {
	// do basic sanity checks
	err := check.ValidateBasic(c.chainID)
	if err != nil {
		return err
	}

	// make sure it has the same validator set we have (static means static)
	if !bytes.Equal(c.Hash(), check.Header.ValidatorsHash) {
		return certerr.ErrValidatorsChanged()
	}

	// then make sure we have the proper signatures for this
	err = c.VSet.VerifyCommit(c.chainID, check.Commit.BlockID,
		check.Header.Height, check.Commit)
	return errors.WithStack(err)
}
