package lightclient

// Broadcaster provides a way to send a signed transaction to a tendermint node
type Broadcaster interface {
	// Broadcast sends into to the chain
	// We only implement BroadcastCommit for now, add others???
	// The return result cannot be fully trusted without downloading signed headers
	Broadcast(tx []byte) (TmBroadcastResult, error)
}

// Checker provides access to calls to get data from the tendermint core
// and all cryptographic proof of its validity
type Checker interface {
	// Prove returns a merkle proof for the given key
	Prove(key []byte) (TmQueryResult, error)

	// SignedHeader gives us Header data along with the backing signatures,
	// so we can validate it externally (matching with the list of
	// known validators)
	SignedHeader(height uint64) (TmSignedHeader, error)

	// WaitForHeight is a useful helper to poll the server until the
	// data is ready for SignedHeader.  Returns nil when the data
	// is present, and error if it aborts.
	WaitForHeight(height uint64) error
}

type Searcher interface {
	// Query gets data from the Blockchain state, possibly with a
	// complex path.  It doesn't worry about proofs
	Query(path string, data []byte) (TmQueryResult, error)
}

// Value represents a database value and is generally a structure
// that can be json serialized.  Bytes() is needed to get the original
// data bytes for validation of proofs
//
// TODO: add Fields() method to get field info???
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
