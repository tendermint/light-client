package certifiers

import (
	"github.com/tendermint/tendermint/types"

	certerr "github.com/tendermint/light-client/certifiers/errors"
)

type Inquiring struct {
	Cert               *Dynamic
	TrustedFullCommits Provider // These are only properly validated data, from local system
	FullCommitSource   Provider // This is a source of new info, like a node rpc, or other import method
}

func NewInquiring(chainID string, fc FullCommit, trusted Provider, source Provider) *Inquiring {
	// store the data in trusted
	trusted.StoreCommit(fc)

	return &Inquiring{
		Cert:               NewDynamic(chainID, fc.Validators, fc.Height()),
		TrustedFullCommits: trusted,
		FullCommitSource:   source,
	}
}

func (c *Inquiring) ChainID() string {
	return c.Cert.ChainID()
}

// Certify makes sure this is checkpoint is valid.
//
// If the validators have changed since the last know time, it looks
// for a path to prove the new validators.
//
// On success, it will store the checkpoint in the store for later viewing
func (c *Inquiring) Certify(commit *Commit) error {
	err := c.useClosestTrust(commit.Height())
	if err != nil {
		return err
	}

	err = c.Cert.Certify(commit)
	if !certerr.IsValidatorsChangedErr(err) {
		return err
	}
	err = c.updateToHash(commit.Header.ValidatorsHash)
	if err != nil {
		return err
	}

	err = c.Cert.Certify(commit)
	if err != nil {
		return err
	}

	// store the new checkpoint
	c.TrustedFullCommits.StoreCommit(
		NewFullCommit(commit, c.Validators()))
	return nil
}

func (c *Inquiring) Validators() *types.ValidatorSet {
	return c.Cert.Cert.vSet
}

func (c *Inquiring) Update(commit *Commit, vals *types.ValidatorSet) error {
	err := c.useClosestTrust(commit.Height())
	if err != nil {
		return err
	}

	err = c.Cert.Update(commit, vals)
	if err == nil {
		c.TrustedFullCommits.StoreCommit(
			NewFullCommit(commit, vals))
	}
	return err
}

func (c *Inquiring) useClosestTrust(h int) error {
	closest, err := c.TrustedFullCommits.GetByHeight(h)
	if err != nil {
		return err
	}

	// if the best seed is not the one we currently use,
	// let's just reset the dynamic validator
	if closest.Height() != c.Cert.LastHeight {
		c.Cert = NewDynamic(c.ChainID(), closest.Validators, closest.Height())
	}
	return nil
}

// updateToHash gets the validator hash we want to update to
// if IsTooMuchChangeErr, we try to find a path by binary search over height
func (c *Inquiring) updateToHash(vhash []byte) error {
	// try to get the match, and update
	fc, err := c.FullCommitSource.GetByHash(vhash)
	if err != nil {
		return err
	}
	err = c.Cert.Update(fc.Commit, fc.Validators)
	// handle IsTooMuchChangeErr by using divide and conquer
	if certerr.IsTooMuchChangeErr(err) {
		err = c.updateToHeight(fc.Height())
	}
	return err
}

// updateToHeight will use divide-and-conquer to find a path to h
func (c *Inquiring) updateToHeight(h int) error {
	// try to update to this height (with checks)
	fc, err := c.FullCommitSource.GetByHeight(h)
	if err != nil {
		return err
	}
	start, end := c.Cert.LastHeight, fc.Height()
	if end <= start {
		return certerr.ErrNoPathFound()
	}
	err = c.Update(fc.Commit, fc.Validators)

	// we can handle IsTooMuchChangeErr specially
	if !certerr.IsTooMuchChangeErr(err) {
		return err
	}

	// try to update to mid
	mid := (start + end) / 2
	err = c.updateToHeight(mid)
	if err != nil {
		return err
	}

	// if we made it to mid, we recurse
	return c.updateToHeight(h)
}
