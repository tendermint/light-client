package util

import (
	"bytes"

	"github.com/pkg/errors"

	lc "github.com/tendermint/light-client"
)

// Auditor takes data, proof, block headers, and tracking of the validator
// set and combines it all to give you complete certainty of the truth
// of a given statement.
type Auditor struct {
	cert lc.Certifier
}

func NewAuditor(cert lc.Certifier) Auditor {
	return Auditor{cert}
}

func (a Auditor) Audit(key, value []byte,
	proof lc.Proof,
	block lc.TmSignedHeader) error {

	root := proof.Root()
	if !proof.Verify(key, value, root) {
		return errors.New("Invalid proof")
	}

	if !bytes.Equal(root, block.Header.AppHash) {
		return errors.New("Header AppHash doesn't match proof")
	}

	if !block.Votes.ForBlock(block.Hash) {
		return errors.New("Votes don't match header")
	}

	// we have traced the data all the way back to a header, now just check
	// this header has all signatures to validate it
	return a.cert.Certify(block)
}
