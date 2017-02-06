package lightclient

import "time"

// Node represents someway to query a tendermint node for info
// Typically via RPC, but could be mocked or connected locally
// TODO: trim this down and distinguish from RPC a bit! (custom types)
type Node interface {
	// Broadcast sends into to the chain
	// The return result cannot be fully trusted without downloading signed headers
	Broadcast(tx []byte) (BroadcastResult, error)

	// Status and Validators give some info, nothing to be trusted though...
	// Unless we find that eg. the ValidatorResult matches the ValidatorHash
	// in a properly signed block header
	Status() (StatusResult, error)
	Validators() (ValidatorResult, error)

	// Query gets data from the Blockchain state, and can optionally provide
	// a Proof we can validate
	Query(path string, data []byte, prove bool) (QueryResult, error)

	// You need to check the Headers and Votes together to prove anything
	// is actually on the chain
	Headers(minHeight, maxHeight int) ([]BlockMeta, error)
	Votes(height int) ([]Vote, error)

	// 	NetInfo() (*ctypes.ResultNetInfo, error)
	// 	DialSeeds(seeds []string) (*ctypes.ResultDialSeeds, error)
	// 	Genesis() (*ctypes.ResultGenesis, error)
	//  Block(height int) (*ctypes.ResultBlock, error)

	// BroadcastTxCommit(tx types.Tx) (*ctypes.ResultBroadcastTx, error)
	// 	BroadcastTxAsync(tx types.Tx) (*ctypes.ResultBroadcastTx, error)
	// 	BroadcastTxSync(tx types.Tx) (*ctypes.ResultBroadcastTx, error)

	// 	ABCIInfo() (*ctypes.ResultABCIInfo, error)
}

type CodeType int32

func (c CodeType) IsOK() bool {
	return int32(c) == 0
}

type BroadcastResult struct {
	Code CodeType `json:"code"` // TODO: rethink this
	Data []byte   `json:"data"`
	Log  string   `json:"log"`
}

func (r BroadcastResult) IsOk() bool {
	return r.Code.IsOK()
}

type StatusResult struct {
	LatestBlockHash   []byte `json:"latest_block_hash"`
	LatestAppHash     []byte `json:"latest_app_hash"`
	LatestBlockHeight int    `json:"latest_block_height"`
	LatestBlockTime   int64  `json:"latest_block_time"` // nano
}

// TODO: how to handle proofs?
// where do we parse them from bytes into Proof objects we can work with
type QueryResult struct {
	Code CodeType `json:"code,omitempty"`
	// Index  int64    `json:"index,omitempty"` // ????
	Key   []byte `json:"key,omitempty"`
	Value []byte `json:"value,omitempty"`
	// Proof Proof  `json:"proof,omitempty"`
	Proof  []byte `json:"proof,omitempty"`
	Height uint64 `json:"height,omitempty"`
	Log    string `json:"log,omitempty"`
}

// Validator more or less from tendermint/types
type Validator struct {
	Address []byte `json:"address"`
	PubKey  []byte `json:"pub_key"`
	// PubKey      crypto.PubKey `json:"pub_key"`
	VotingPower int64 `json:"voting_power"`
}

type ValidatorResult struct {
	BlockHeight int
	Validators  []Validator
}

type BlockMeta struct {
	Hash   []byte
	Header Header
}

// Header is the info in block headers (from tendermint/types/block.go)
type Header struct {
	ChainID        string    `json:"chain_id"`
	Height         int       `json:"height"`
	Time           time.Time `json:"time"`
	NumTxs         int       `json:"num_txs"` // XXX: Can we get rid of this?
	LastBlockID    []byte    `json:"last_block_id"`
	LastCommitHash []byte    `json:"last_commit_hash"` // commit from validators from the last block
	DataHash       []byte    `json:"data_hash"`        // transactions
	ValidatorsHash []byte    `json:"validators_hash"`  // validators for the current block
	AppHash        []byte    `json:"app_hash"`         // state after txs from the previous block
}

// Vote must be verified by the Node implementation, this asserts a validly signed
// precommit vote for the given Height and BlockHash.
// The client can decide if these validators are to be trusted.
type Vote struct {
	ValidatorAddress []byte `json:"validator_address"`
	// ValidatorIndex   int              `json:"validator_index"`
	Height    int    `json:"height"`
	BlockHash []byte `json:"block_hash"`
	// Round            int              `json:"round"`
	// Type             byte             `json:"type"`
	// BlockID          BlockID          `json:"block_id"` // zero if vote is nil.
	// Signature        crypto.Signature `json:"signature"`
}
