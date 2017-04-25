package lightclient

import (
	"bytes"

	"github.com/pkg/errors"
	rtypes "github.com/tendermint/tendermint/rpc/core/types"
	"github.com/tendermint/tendermint/types"
)

// Certifier checks the votes to make sure the block really is signed properly.
// Certifier must know the current set of validitors by some other means.
type Certifier interface {
	Certify(check Checkpoint) error
}

// Checkpoint is basically the rpc /commit response, but extended
//
// This is the basepoint for proving anything on the blockchain. It contains
// a signed header.  If the signatures are valid and > 2/3 of the known set,
// we can store this checkpoint and use it to prove any number of aspects of
// the system: such as txs, abci state, validator sets, etc...
type Checkpoint struct {
	Header *types.Header `json:"header"`
	Commit *types.Commit `json:"commit"`
}

func CheckpointFromResult(commit *rtypes.ResultCommit) Checkpoint {
	return Checkpoint{
		Header: commit.Header,
		Commit: commit.Commit,
	}
}

func (c Checkpoint) Height() int {
	if c.Header == nil {
		return 0
	}
	return c.Header.Height
}

// ValidateBasic does basic consistency checks and makes sure the headers
// and commits are all consistent and refer to our chain.
//
// Make sure to use a Verifier to validate the signatures actually provide
// a significantly strong proof for this header's validity.
func (c Checkpoint) ValidateBasic(chainID string) error {
	// make sure the header is reasonable
	if c.Header == nil {
		return errors.New("Checkpoint missing header")
	}
	if c.Header.ChainID != chainID {
		return errors.Errorf("Header belongs to another chain '%s' not '%s'",
			c.Header.ChainID, chainID)
	}

	if c.Commit == nil {
		return errors.New("Checkpoint missing commits")
	}

	// make sure the header and commit match (height and hash)
	if c.Commit.Height() != c.Header.Height {
		return ErrHeightMismatch(c.Commit.Height(), c.Header.Height)
	}
	hhash := c.Header.Hash()
	chash := c.Commit.BlockID.Hash
	if !bytes.Equal(hhash, chash) {
		return errors.Errorf("Commits sign block %X header is block %X",
			chash, hhash)
	}

	// make sure the commit is reasonable
	err := c.Commit.ValidateBasic()
	if err != nil {
		return errors.WithStack(err)
	}

	// looks good, we just need to make sure the signatures are really from
	// empowered validators
	return nil
}
