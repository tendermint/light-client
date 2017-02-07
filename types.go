package lightclient

import (
	"bytes"
	"time"
)

// TmCode is a code from Tendermint
type TmCode int32

func (c TmCode) IsOK() bool {
	return int32(c) == 0
}

type TmBroadcastResult struct {
	Code TmCode `json:"code"` // TODO: rethink this
	Data []byte `json:"data"`
	Log  string `json:"log"`
}

func (r TmBroadcastResult) IsOk() bool {
	return r.Code.IsOK()
}

type TmStatusResult struct {
	LatestBlockHash   []byte `json:"latest_block_hash"`
	LatestAppHash     []byte `json:"latest_app_hash"`
	LatestBlockHeight int    `json:"latest_block_height"`
	LatestBlockTime   int64  `json:"latest_block_time"` // nano
}

// TODO: how to handle proofs?
// where do we parse them from bytes into Proof objects we can work with
type TmQueryResult struct {
	Code TmCode `json:"code"`
	// Index  int64    `json:"index,omitempty"` // ????
	Key   []byte `json:"key"`
	Value []byte `json:"value"`
	// Proof Proof  `json:"proof"`
	Proof  []byte `json:"proof"`
	Height uint64 `json:"height"`
	Log    string `json:"log"`
}

// TmValidator more or less from tendermint/types
type TmValidator struct {
	Address []byte `json:"address"`
	PubKey  []byte `json:"pub_key"`
	// PubKey      crypto.PubKey `json:"pub_key"`
	VotingPower int64 `json:"voting_power"`
}

type TmValidatorResult struct {
	BlockHeight int
	Validators  []TmValidator
}

// TmBlockMeta is the Header info and the Hash that corresponds to it
// (and which is used to cannonically identiry the block)
// The Node implementation is responsible for validating this is correct,
// thus we can return the Header is any useful format, not byte-for-byte how
// tendermint stores it.
type TmBlockMeta struct {
	Hash   []byte
	Header TmHeader
}

// TmHeader is the info in block headers (from tendermint/types/block.go)
type TmHeader struct {
	ChainID        string    `json:"chain_id"`
	Height         int       `json:"height"`
	Time           time.Time `json:"time"`    // or int64 nanoseconds????
	NumTxs         int       `json:"num_txs"` // XXX: Can we get rid of this?
	LastBlockID    []byte    `json:"last_block_id"`
	LastCommitHash []byte    `json:"last_commit_hash"` // commit from validators from the last block
	DataHash       []byte    `json:"data_hash"`        // transactions
	ValidatorsHash []byte    `json:"validators_hash"`  // validators for the current block
	AppHash        []byte    `json:"app_hash"`         // state after txs from the previous block
}

// TmVote must be verified by the Node implementation, this asserts a validly
// signed precommit vote for the given Height and BlockHash.
// The client can decide if these validators are to be trusted.
type TmVote struct {
	ValidatorAddress []byte `json:"validator_address"`
	// ValidatorIndex   int              `json:"validator_index"`
	Height    int    `json:"height"`
	BlockHash []byte `json:"block_hash"`
	// Round            int              `json:"round"`
	// Type             byte             `json:"type"`
	// BlockID          BlockID          `json:"block_id"` // zero if vote is nil.
	// Signature        crypto.Signature `json:"signature"`
}

// TmVotes is a slice of TmVote structs, but let's add some control access here
type TmVotes []TmVote

// ForBlock returns true only if all votes are for the given block
func (v TmVotes) ForBlock(hash []byte) bool {
	if len(v) == 0 {
		return false
	}

	for _, vv := range v {
		if !bytes.Equal(hash, vv.BlockHash) {
			return false
		}
	}

	return true
}
