package lightclient

import (
	"bytes"

	"github.com/pkg/errors"
	rtypes "github.com/tendermint/tendermint/rpc/core/types"
	"github.com/tendermint/tendermint/types"
)

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

func NewCheckpoint(commit *rtypes.ResultCommit) Checkpoint {
	return Checkpoint{
		Header: commit.Header,
		Commit: commit.Commit,
	}
}

func (c Checkpoint) Height() int {
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

	// make sure the commit is reasonable
	if c.Commit == nil {
		return errors.New("Checkpoint missing commits")
	}
	err := c.Commit.ValidateBasic()
	if err != nil {
		return errors.WithStack(err)
	}

	// make sure the header and commit match (height and hash)
	if c.Commit.Height() != c.Header.Height {
		return errors.Errorf("Commit height %d != header height %d",
			c.Commit.Height(), c.Header.Height)
	}
	hhash := c.Header.Hash()
	chash := c.Commit.BlockID.Hash
	if !bytes.Equal(hhash, chash) {
		return errors.Errorf("Commits sign block %X header is block %X",
			chash, hhash)
	}

	// looks good, we just need to make sure the signatures are really from
	// empowered validators
	return nil
}

// CheckValidators should only be used after you fully trust this checkpoint
//
// It checks if these really are the validators authorized to sign the
// checkpoint.
func (c Checkpoint) CheckValidators(vals []*types.Validator) error {
	if len(vals) == 0 {
		return errors.New("No validators provided")
	}
	hash := types.NewValidatorSet(vals).Hash()
	if !bytes.Equal(hash, c.Header.ValidatorsHash) {
		return errors.New("Validator hashes differ")
	}
	return nil
}

// CheckTxs checks if the entire set of transactions for the block matches
// the Checkpoint header.
func (c Checkpoint) CheckTxs(txs types.Txs) error {
	hash := txs.Hash()
	if !bytes.Equal(hash, c.Header.DataHash) {
		return errors.Errorf("TxHash %X != Header hash %X",
			hash, c.Header.DataHash)
	}
	return nil
}

// TODO: one tx plus proof.... need changes in the
// func (c Checkpoint) CheckTx(tx types.Tx) error {
//   return nil
// }

// CheckAppState validates whether the key-value pair and merkle proof
// can be verified with this Checkpoint.
func (c Checkpoint) CheckAppState(k, v []byte, proof Proof) error {
	if !bytes.Equal(proof.Root(), c.Header.AppHash) {
		return errors.Errorf("Proof is for AppHash %X but header has %X",
			proof.Root(), c.Header.AppHash)
	}
	if !proof.Verify(k, v, c.Header.AppHash) {
		return errors.New("Proof doesn't match given key-value pair")
	}
	return nil
}
