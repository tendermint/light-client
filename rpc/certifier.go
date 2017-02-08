package rpc

import (
	"bytes"

	"github.com/pkg/errors"
	lc "github.com/tendermint/light-client"
)

// StaticCertifier assumes a static set of validators, set on
// initilization and checks against them.
//
// Good for testing or really simple chains.  You will want a
// better implementation when the validator set can actually change.
type StaticCertifier struct {
	vals []lc.TmValidator
}

func (c StaticCertifier) assertCertifier() lc.Certifier {
	return c
}

func (c StaticCertifier) Certify(block lc.TmSignedHeader) error {
	var points int64
	// check the validator sigs are valid
	for _, v := range block.Votes {
		val := c.valByAddr(v.ValidatorAddress)
		if val == nil {
			continue // ignore unknown validators
		}
		valid := val.PubKey.VerifyBytes(v.SignBytes, v.Signature)
		if !valid {
			return errors.Errorf("Invalid signature for validator: %X",
				v.ValidatorAddress)
		}
		points += val.VotingPower
	}

	// make sure there were enough
	total := c.totalVotes()
	if 3*points <= 2*total {
		return errors.Errorf("Only %d out of %d votes for this block",
			points, total)
	}

	return nil
}

// valByAddr looks up the validator (and public key) by the known
// address.
// FIXME: inefficient for large sets, but this is not for that
func (c StaticCertifier) valByAddr(addr []byte) *lc.TmValidator {
	for i := range c.vals {
		if bytes.Equal(addr, c.vals[i].Address) {
			return &c.vals[i]
		}
	}
	return nil
}

// totalVotes returns the total voting power
// FIXME: could be cached for efficiency, but key is simplicity
func (c StaticCertifier) totalVotes() int64 {
	var votes int64
	for i := range c.vals {
		votes += c.vals[i].VotingPower
	}
	return votes
}
