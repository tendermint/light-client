package certifiers

import (
	"fmt"

	lc "github.com/tendermint/light-client"
	"github.com/tendermint/tendermint/types"
)

var _ lc.Certifier = &DynamicCertifier{}

// DynamicCertifier uses a StaticCertifier to evaluate the checkpoint
// but allows for a change, if we present enough proof
//
// TODO: do we keep a long history so we can use our memory to validate
// checkpoints from previously valid validator sets????
type DynamicCertifier struct {
	Cert       *StaticCertifier
	LastHeight int
}

func NewDynamic(chainID string, vals *types.ValidatorSet) *DynamicCertifier {
	return &DynamicCertifier{
		Cert:       NewStatic(chainID, vals),
		LastHeight: 0,
	}
}

// Certify handles this with
func (c *DynamicCertifier) Certify(check lc.Checkpoint) error {
	err := c.Cert.Certify(check)
	if err == nil {
		// update last seen height if input is valid
		c.LastHeight = check.Height()
	}
	return err
}

// Update will verify if this is a valid change and update
// the certifying validator set if safe to do so.
//
// Returns an error if update is impossible (invalid proof or IsTooMuchChangeErr)
func (c *DynamicCertifier) Update(check lc.Checkpoint, vset *types.ValidatorSet) error {
	// ignore all checkpoints in the past -> only to the future
	if check.Height() <= c.LastHeight {
		return ErrPastTime()
	}

	// first, verify if the input is self-consistent....
	err := check.ValidateBasic(c.Cert.ChainID)
	if err != nil {
		return err
	}

	// TODO: now, make sure not too much change... meaning this commit
	// would be approved by the currently known validator set
	// as well as the new set
	err = VerifyCommitAny(c.Cert.VSet, vset, c.Cert.ChainID,
		check.Commit.BlockID, check.Header.Height, check.Commit)
	if err != nil {
		return ErrTooMuchChange()
	}

	// looks good, we can update
	c.Cert = NewStatic(c.Cert.ChainID, vset)
	c.LastHeight = check.Height()
	return nil
}

// VerifyCommitAny will check to see if the set would
// be valid with a different validator set.
//
// old is the validator set that we know
// * over 2/3 of the power in old signed this block
//
// cur is the validator set that signed this block
// * only votes from old are sufficient for 2/3 majority
//   in the new set as well
//
// That means that:
// * 10% of the valset can't just declare themselves kings
// * If the validator set is 3x old size, we need more proof to trust
//
// *** TODO: move this.
// It belongs in tendermint/types/validator_set.go: VerifyCommitAny
func VerifyCommitAny(old, cur *types.ValidatorSet, chainID string,
	blockID types.BlockID, height int, commit *types.Commit) error {

	if cur.Size() != len(commit.Precommits) {
		return fmt.Errorf("Invalid commit -- wrong set size: %v vs %v", cur.Size(), len(commit.Precommits))
	}
	if height != commit.Height() {
		return fmt.Errorf("Invalid commit -- wrong height: %v vs %v", height, commit.Height())
	}

	oldVotingPower := int64(0)
	curVotingPower := int64(0)
	seen := map[int]bool{}
	round := commit.Round()

	for idx, precommit := range commit.Precommits {
		// first check as in VerifyCommit
		if precommit == nil {
			continue
		}
		if precommit.Height != height {
			return lc.ErrHeightMismatch(height, precommit.Height)
		}
		if precommit.Round != round {
			return fmt.Errorf("Invalid commit -- wrong round: %v vs %v", round, precommit.Round)
		}
		if precommit.Type != types.VoteTypePrecommit {
			return fmt.Errorf("Invalid commit -- not precommit @ index %v", idx)
		}
		if !blockID.Equals(precommit.BlockID) {
			continue // Not an error, but doesn't count
		}

		// we only grab by address, ignoring unknown validators
		vi, ov := old.GetByAddress(precommit.ValidatorAddress)
		if ov == nil || seen[vi] {
			continue // missing or double vote...
		}
		seen[vi] = true

		// Validate signature old school
		precommitSignBytes := types.SignBytes(chainID, precommit)
		if !ov.PubKey.VerifyBytes(precommitSignBytes, precommit.Signature) {
			return fmt.Errorf("Invalid commit -- invalid signature: %v", precommit)
		}
		// Good precommit!
		oldVotingPower += ov.VotingPower

		// check new school
		_, cv := cur.GetByIndex(idx)
		if cv.PubKey.Equals(ov.PubKey) {
			// make sure this is properly set in the current block as well
			curVotingPower += cv.VotingPower
		}
	}

	if oldVotingPower <= old.TotalVotingPower()*2/3 {
		return fmt.Errorf("Invalid commit -- insufficient old voting power: got %v, needed %v",
			oldVotingPower, (old.TotalVotingPower()*2/3 + 1))
	} else if curVotingPower <= cur.TotalVotingPower()*2/3 {
		return fmt.Errorf("Invalid commit -- insufficient cur voting power: got %v, needed %v",
			curVotingPower, (cur.TotalVotingPower()*2/3 + 1))
	}
	return nil
}
