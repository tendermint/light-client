package lightclient

import (
	"bytes"
	"time"

	crypto "github.com/tendermint/go-crypto"
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

// TmQueryResult stores all info from a query or a proof
// The Checker/Searcher is responsible for parsing the Value and Proof into
// usable formats (not just opaque bytes) via their Reader
type TmQueryResult struct {
	Code   TmCode `json:"code"`
	Key    []byte `json:"key"`
	Value  Value  `json:"value"`
	Proof  Proof  `json:"proof"`
	Height uint64 `json:"height"`
	Log    string `json:"log"`
}

// TmSignedHeader returns the header info and the corresponding validator
// signatures for it (even if those are on differrent rpc endpoints).
//
// Checker is responsible for validating that the Hash matches the Header
// which lets us return the Header in a non-binary format.
// Votes must also be pre-verified, signatures checked, etc.
type TmSignedHeader struct {
	Hash       []byte
	Header     TmHeader // contains height
	Votes      TmVotes
	LastHeight uint64 // the most recent block commited to the chain
}

// TmHeader is the info in block headers (from tendermint/types/block.go)
type TmHeader struct {
	ChainID string    `json:"chain_id"`
	Height  uint64    `json:"height"`
	Time    time.Time `json:"time"` // or int64 nanoseconds????
	// NumTxs         int       `json:"num_txs"` // XXX: Can we get rid of this?
	LastBlockID    []byte `json:"last_block_id"`
	LastCommitHash []byte `json:"last_commit_hash"` // commit from validators from the last block
	DataHash       []byte `json:"data_hash"`        // transactions
	ValidatorsHash []byte `json:"validators_hash"`  // validators for the current block
	AppHash        []byte `json:"app_hash"`         // state after txs from the previous block
}

// TmVote must be verified by the Node implementation, this asserts a validly
// signed precommit vote for the given Height and BlockHash.
// The client can decide if these validators are to be trusted.
type TmVote struct {
	// SignBytes are the cannonical bytes the signature refers to
	// This is verified in the Checker
	SignBytes []byte `json:"sign_bytes"`

	// Signature and ValidatorAddress represent who signed this
	// This information is not verified and must be validated
	// by the caller
	Signature        crypto.Signature `json:"signature"`
	ValidatorAddress []byte           `json:"validator_address"`

	// Height and BlockHash is embedded in SignBytes
	Height    uint64 `json:"height"`
	BlockHash []byte `json:"block_hash"`
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

type TmStatusResult struct {
	LatestBlockHash   []byte `json:"latest_block_hash"`
	LatestAppHash     []byte `json:"latest_app_hash"`
	LatestBlockHeight int    `json:"latest_block_height"`
	LatestBlockTime   int64  `json:"latest_block_time"` // nano
}

// TmValidator more or less from tendermint/types
type TmValidator struct {
	Address []byte `json:"address"`
	// PubKey  []byte `json:"pub_key"`
	PubKey      crypto.PubKey `json:"pub_key"`
	VotingPower int64         `json:"voting_power"`
}

type TmValidatorResult struct {
	BlockHeight uint64
	Validators  []TmValidator
}
