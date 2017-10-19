package certifiers

import (
	"github.com/tendermint/tendermint/types"
)

type InquiringCertifier struct {
	Cert         *DynamicCertifier
	TrustedSeeds Provider // These are only properly validated data, from local system
	SeedSource   Provider // This is a source of new info, like a node rpc, or other import method
}

func NewInquiring(chainID string, seed Seed, trusted Provider, source Provider) *InquiringCertifier {
	// store the data in trusted
	trusted.StoreSeed(seed)

	return &InquiringCertifier{
		Cert:         NewDynamic(chainID, seed.Validators, seed.Height()),
		TrustedSeeds: trusted,
		SeedSource:   source,
	}
}

func (c *InquiringCertifier) ChainID() string {
	return c.Cert.Cert.ChainID
}

// Certify makes sure this is checkpoint is valid.
//
// If the validators have changed since the last know time, it looks
// for a path to prove the new validators.
//
// On success, it will store the checkpoint in the store for later viewing
func (c *InquiringCertifier) Certify(check Checkpoint) error {
	err := c.useClosestTrust(check.Height())
	if err != nil {
		return err
	}

	err = c.Cert.Certify(check)
	if !IsValidatorsChangedErr(err) {
		return err
	}
	err = c.updateToHash(check.Header.ValidatorsHash)
	if err != nil {
		return err
	}

	err = c.Cert.Certify(check)
	if err != nil {
		return err
	}

	// store the new checkpoint
	c.TrustedSeeds.StoreSeed(Seed{
		Checkpoint: check,
		Validators: c.Cert.Cert.VSet,
	})
	return nil
}

func (c *InquiringCertifier) Update(check Checkpoint, vals *types.ValidatorSet) error {
	err := c.useClosestTrust(check.Height())
	if err != nil {
		return err
	}

	err = c.Cert.Update(check, vals)
	if err == nil {
		c.TrustedSeeds.StoreSeed(Seed{Checkpoint: check, Validators: vals})
	}
	return err
}

func (c *InquiringCertifier) useClosestTrust(h int) error {
	closest, err := c.TrustedSeeds.GetByHeight(h)
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
func (c *InquiringCertifier) updateToHash(vhash []byte) error {
	// try to get the match, and update
	seed, err := c.SeedSource.GetByHash(vhash)
	if err != nil {
		return err
	}
	err = c.Cert.Update(seed.Checkpoint, seed.Validators)
	// handle IsTooMuchChangeErr by using divide and conquer
	if IsTooMuchChangeErr(err) {
		err = c.updateToHeight(seed.Height())
	}
	return err
}

// updateToHeight will use divide-and-conquer to find a path to h
func (c *InquiringCertifier) updateToHeight(h int) error {
	// try to update to this height (with checks)
	seed, err := c.SeedSource.GetByHeight(h)
	if err != nil {
		return err
	}
	start, end := c.Cert.LastHeight, seed.Height()
	if end <= start {
		return ErrNoPathFound()
	}
	err = c.Update(seed.Checkpoint, seed.Validators)

	// we can handle IsTooMuchChangeErr specially
	if !IsTooMuchChangeErr(err) {
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
