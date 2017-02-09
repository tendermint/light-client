package lightclient

// Proof is a generalization of merkle.IAVLProof and represents any
// merkle proof that can validate a key-value pair back to a root hash.
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
