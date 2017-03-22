package lightclient

import (
	keys "github.com/tendermint/go-keys"
	merkle "github.com/tendermint/go-merkle"
)

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

// ProofReader is an abstraction to let us parse proofs
type ProofReader interface {
	ReadProof(data []byte) (Proof, error)
}

// MerkleReader is currently the only implementation of ProofReader,
// using the IAVLProof from go-merkle
var MerkleReader ProofReader = merkleReader{}

type merkleReader struct{}

func (p merkleReader) ReadProof(data []byte) (Proof, error) {
	return merkle.ReadProof(data)
}

type Value interface {
	Bytes() []byte
}

// ValueReader is an abstraction to let us parse application-specific values
type ValueReader interface {
	// ReadValue accepts a key, value pair to decode.  The value bytes must be
	// retained in the returned Value implementation.
	//
	// key *may* be present and can be used as a hint of how to parse the data
	// when your application handles multiple formats
	ReadValue(key, value []byte) (Value, error)
}

type SignableReader interface {
	ReadSignable(data []byte) (keys.Signable, error)
}
