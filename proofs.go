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
	Certify(votes Votes, height int) error
}

// Auditor takes data, proof, block headers, and tracking of the validator
// set and combines it all to give you complete certainty of the truth
// of a given statement.
type Auditor interface {
	Audit(key, value []byte, proof Proof, header BlockMeta, votes Votes) error
}

// TODO: some glue code to query proof and header, and validate all aspects
func NewAuditor(cert Certifier) Auditor {
	return auditor{cert}
}

type auditor struct {
	cert Certifier
}

func (a auditor) Audit(key, value []byte,
	proof Proof,
	header BlockMeta,
	votes Votes) error {

	root := proof.Root()
	if !proof.Verify(key, value, root) {
		return errors.New("Invalid proof")
	}

	if !bytes.Equal(root, header.Header.AppHash) {
		return errors.New("Header AppHash doesn't match proof")
	}

	if !votes.ForBlock(header.Hash) {
		return errors.New("Votes don't match header")
	}

	// we have traced the data all the way back to a header, now just check
	// this header has all signatures to validate it
	return a.cert.Certify(votes, header.Header.Height)
}
