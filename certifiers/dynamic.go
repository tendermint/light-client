package certifiers

import (
	rawerr "errors"

	"github.com/pkg/errors"
	lc "github.com/tendermint/light-client"
	"github.com/tendermint/tendermint/types"
)

var (
	errTooMuchChange = rawerr.New("Validators differ between header and certifier")
	errPastTime      = rawerr.New("Update older than certifier height")
)

// ToMuchChange asserts whether and error is due to too much change
// between these validators sets
func ToMuchChange(err error) bool {
	return errors.Cause(err) == errTooMuchChange
}

// DynamicCertifier uses a StaticCertifier to evaluate the checkpoint
// but allows for a change, if we present enough proof
//
// TODO: do we keep a long history so we can use our memory to validate
// checkpoints from previously valid validator sets????
type DynamicCertifier struct {
	Cert       *StaticCertifier
	LastHeight int
}

func NewDynamic(chainID string, vals []*types.Validator) *DynamicCertifier {
	return &DynamicCertifier{
		Cert:       NewStatic(chainID, vals),
		LastHeight: 0,
	}
}

func (c *DynamicCertifier) assertCertifier() lc.Certifier {
	return c
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
// Returns an error if update is impossible (invalid proof or TooMuchChange)
func (c *DynamicCertifier) Update(check lc.Checkpoint, vals []*types.Validator) error {
	// ignore all checkpoints in the past -> only to the future
	if check.Height() <= c.LastHeight {
		return errors.WithStack(errPastTime)
	}

	// first, verify if the input is self-consistent....
	st := NewStatic(c.Cert.ChainID, vals)
	err := st.Certify(check)
	if err != nil {
		return err
	}

	// TODO: now, make sure not too much change

	// looks good, we can update
	c.Cert = st
	c.LastHeight = check.Height()
	return nil
}
