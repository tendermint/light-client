package lightclient

// TODO: how do we parse this?
// Hard-code the use of merkle.IAVLProof???
// Some more clever way?
type Proof interface {
	// Root returns the RootHash of the merkle tree used in the proof,
	// This is important for correlating it with a block header.
	Root() []byte

	// Verify returns true iff this proof validates this key and value belong
	// to the given root
	Verify(key, value, root []byte) bool
}

// TODO: some interface to track the validator set (various algorithms)

// TODO: some glue code to query proof and header, and validate all aspects

// TODO: someway to save/export a given proof for another client??
