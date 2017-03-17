package certifiers

import (
	lc "github.com/tendermint/light-client"
	"github.com/tendermint/tendermint/types"
)

type InquiringCertifier struct {
	Cert *DynamicCertifier
	Provider
}

func (c *InquiringCertifier) Certify(check lc.Checkpoint) error {
	err := c.Cert.Certify(check)
	if !ValidatorsChanged(err) {
		return err
	}
	err = c.updateToHash(check.Height(), check.Header.ValidatorsHash)
	if err != nil {
		return err
	}
	return c.Cert.Certify(check)
}

func (c *InquiringCertifier) Update(check lc.Checkpoint, vals []*types.Validator) error {
	err := c.Cert.Update(check, vals)
	if err == nil {
		c.StoreSeed(Seed{Checkpoint: check, Validators: vals})
	}
	return err
}

func (c *InquiringCertifier) updateToHash(h int, vhash []byte) error {
	// try to get the match, and update
	seed, err := c.GetByHash(vhash)
	if err != nil {
		return err
	}
	err = c.Cert.Update(seed.Checkpoint, seed.Validators)
	// TODO: handle TooMuchChange by using divide and conquer
	return err
}
