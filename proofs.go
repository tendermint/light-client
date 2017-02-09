package lightclient

import (
	"bytes"

	"github.com/pkg/errors"
)

// TODO: how do we parse this?
// Hard-code the use of merkle.IAVLProof???
// Some more clever way?
// TODO: someway to save/export a given proof for another client??
type Proof interface {
	// Root returns the RootHash of the merkle tree used in the proof,
	// This is important for correlating it with a block header.
	Root() []byte

	// Verify returns true iff this proof validates this key and value belong
	// to the given root
	Verify(key, value, root []byte) bool
}

// Certifier checks the votes to make sure the block really is signed properly.
// Certifier must know the current set of validitors by some other means.
// TODO: some implementation to track the validator set (various algorithms)
type Certifier interface {
	Certify(block TmSignedHeader) error
}

// Auditor takes data, proof, block headers, and tracking of the validator
// set and combines it all to give you complete certainty of the truth
// of a given statement.
//
// TODO: move this into some sort of util package that does mashups based
// solely on interface types
type Auditor struct {
	cert Certifier
}

func NewAuditor(cert Certifier) Auditor {
	return Auditor{cert}
}

func (a Auditor) Audit(key, value []byte,
	proof Proof,
	block TmSignedHeader) error {

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
