package lightclient

// Prover is anything that can provide proofs.
// Such as a AppProver (for merkle proofs of app state)
// or TxProver (for merkle proofs that a tx is in a block)
type Prover interface {
	// Get returns the key for the given block height
	// The prover should accept h=0 for latest height
	Get(key []byte, h uint64) (Proof, error)
	Unmarshal([]byte) (Proof, error)
}

// Proof is a generic interface for data along with the cryptographic proof
// of it's validity, tied to a checkpoint.
//
// Every implementation should offer some method to recover the data itself
// that was proven (like k-v pair, tx bytes, etc....)
type Proof interface {
	BlockHeight() uint64
	Validate(Checkpoint) error // Make sure the checkpoint is validated and proper height
	// Marshal prepares for storage
	Marshal() ([]byte, error)
	// Data extracts the query result we want to see
	Data() []byte
}
