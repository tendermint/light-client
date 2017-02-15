package rpc

import (
	"bytes"
	"time"

	"github.com/pkg/errors"
	merkle "github.com/tendermint/go-merkle"
	wire "github.com/tendermint/go-wire"
	lc "github.com/tendermint/light-client"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	ttypes "github.com/tendermint/tendermint/types"
)

type Node struct {
	client *HTTPClient
	// this is needed to calculate sign bytes for votes
	chainID string
	lc.ProofReader
	lc.ValueReader
}

// MerkleReader is currently the only implementation of ProofReader,
// using the IAVLProof from go-merkle
var MerkleReader lc.ProofReader = merkleReader{}

type merkleReader struct{}

func (p merkleReader) ReadProof(data []byte) (lc.Proof, error) {
	return merkle.ReadProof(data)
}

func NewNode(rpcAddr, chainID string, valuer lc.ValueReader) Node {
	return Node{
		client:      NewClient(rpcAddr, "/websocket"),
		chainID:     chainID,
		ProofReader: MerkleReader,
		ValueReader: valuer,
	}
}

func (n Node) assertBroadcaster() lc.Broadcaster {
	return n
}

func (n Node) assertChecker() lc.Checker {
	return n
}

func (n Node) assertSearcher() lc.Searcher {
	return n
}

// Broadcast sends the transaction to a tendermint node and waits until
// it is committed.
//
// If it failed on CheckTx, we return the result of CheckTx, otherwise
// we return the result of DeliverTx
func (n Node) Broadcast(tx []byte) (res lc.TmBroadcastResult, err error) {
	cr, err := n.client.BroadcastTxCommit(tx)
	if err != nil {
		return
	}
	if cr.DeliverTx != nil {
		d := cr.DeliverTx
		res.Code = lc.TmCode(d.Code)
		res.Data = d.Data
		res.Log = d.Log
	} else {
		c := cr.CheckTx
		res.Code = lc.TmCode(c.Code)
		res.Data = c.Data
		res.Log = c.Log
	}
	return
}

// Query gets data from the Blockchain state, possibly with a
// complex path.  It doesn't worry about proofs
func (n Node) Query(path string, data []byte) (lc.TmQueryResult, error) {
	qr, err := n.client.ABCIQuery(path, data, false)
	return n.queryResp(qr, err)
}

// Prove returns a merkle proof for the given key
func (n Node) Prove(key []byte) (lc.TmQueryResult, error) {
	qr, err := n.client.ABCIQuery("/key", key, true)
	return n.queryResp(qr, err)
}

func (n Node) queryResp(qr *ctypes.ResultABCIQuery, err error) (lc.TmQueryResult, error) {
	if qr == nil || err != nil {
		return lc.TmQueryResult{}, err
	}
	r := qr.Response
	res := lc.TmQueryResult{
		Code:   lc.TmCode(r.Code),
		Key:    r.Key,
		Log:    r.Log,
		Height: r.Height,
	}
	// parse the value if it exists
	if len(r.Value) > 0 {
		res.Value, err = n.ReadValue(r.Key, r.Value)
	}
	// parse the proof if it exists
	if len(r.Proof) > 0 {
		res.Proof, err = n.ReadProof(r.Proof)
	}
	return res, err
}

// SignedHeader gives us Header data along with the backing signatures
//
// It is also responsible for making the blockhash and header match,
// and that all votes are valid pre-commit votes for this block
// It does not check if the keys signing the votes are actual validators
func (n Node) SignedHeader(height uint64) (lc.TmSignedHeader, error) {
	h := int(height)

	// get data from rpc
	ci, err := n.getCommitInfo(h)
	if err != nil {
		return lc.TmSignedHeader{}, err
	}

	// validate and process it
	res, err := n.validateCommitInfo(ci)
	if err != nil {
		return lc.TmSignedHeader{}, err
	}

	// make sure the height is what we wanted
	if res.Height() != height {
		err = errors.Errorf("Returned header for height %d, not %d",
			res.Height(), h)
	}
	return res, err
}

// ExportSignedHeader downloads and verifies the same info as
// SignedHeader, but returns a serialized version of the proof to
// be stored for later use.
//
// The result should be considered opaque bytes, but can be passed into
// ImportSignedHeader to get data ready for a Certifier
func (n Node) ExportSignedHeader(height uint64) ([]byte, error) {
	h := int(height)

	// get data from rpc
	ci, err := n.getCommitInfo(h)
	if err != nil {
		return nil, err
	}

	// validate and process it
	res, err := n.validateCommitInfo(ci)
	if err != nil {
		return nil, err
	}

	// make sure the height is what we wanted
	if res.Height() != height {
		err = errors.Errorf("Returned header for height %d, not %d",
			res.Height(), h)
	}

	// serialize data for later use
	return ci.Bytes(), err
}

// ImportSignedHeader takes serialized data from ExportSignedHeader
// and verifies and processes it the same as SignedHeader.
//
// The result can be used just as the result from SignedHeader, and
// passed to a Certifier
func (n Node) ImportSignedHeader(data []byte) (lc.TmSignedHeader, error) {
	ci, err := loadCommitInfo(data)
	if err != nil {
		return lc.TmSignedHeader{}, err
	}
	// validate and process it
	return n.validateCommitInfo(ci)
}

// Wait for height will poll status at reasonable intervals until
// we can safely call SignedHeader at the given block height.
// This means that both the block header itself, as well as all
// validator signatures are available
//
// In this current implementation, we must wait until height+1,
// as the signatures are in the following block.
func (n Node) WaitForHeight(height uint64) error {
	h := int(height) + 1 // off-by-one shizzit
	wait := 1
	for wait > 0 {
		s, err := n.client.Status()
		if err != nil {
			return err
		}
		wait = h - s.LatestBlockHeight
		if wait > 10 {
			return errors.Errorf("Waiting for %d block... aborting", wait)
		} else if wait > 0 {
			// estimate of wait time....
			// wait half a second for the next block (in progress)
			// plus one second for every full block
			delay := time.Duration(wait-1)*time.Second + 500*time.Millisecond
			time.Sleep(delay)
		}
	}
	// guess we waited long enough
	return nil
}

// UNSAFE - use wisely
//
// Validators returns the current set of validators from the
// node we call.  There is no guarantee it is correct.
//
// This is intended for use in test cases, or to query many nodes
// to find consensus before trusting it.
func (n Node) Validators() (lc.TmValidatorResult, error) {
	vres, err := n.client.Validators()
	if err != nil {
		return lc.TmValidatorResult{}, err
	}
	// now we transform them into our world
	vals := vres.Validators
	rvals := make([]lc.TmValidator, len(vals))

	res := lc.TmValidatorResult{
		BlockHeight: uint64(vres.BlockHeight),
		Validators:  rvals,
	}
	for i, v := range vals {
		rvals[i] = lc.TmValidator{
			Address:     v.Address,
			VotingPower: v.VotingPower,
			PubKey:      v.PubKey,
		}
	}
	return res, nil
}

type commitInfo struct {
	Header    *ttypes.BlockMeta
	MaxHeight int
	Commit    *ttypes.Commit
}

func (c commitInfo) Bytes() []byte {
	return wire.BinaryBytes(c)
}

func loadCommitInfo(data []byte) (res commitInfo, err error) {
	err = wire.ReadBinaryBytes(data, &res)
	return
}

func (n Node) getCommitInfo(h int) (res commitInfo, err error) {
	// we get the raw data first...
	res.Header, res.MaxHeight, err = n.getHeader(h)
	if err != nil {
		return
	}
	res.Commit, err = n.getCommit(h)
	return
}

// get header returns the header info along with the most recent height
func (n Node) getHeader(h int) (*ttypes.BlockMeta, int, error) {
	bis, err := n.client.BlockchainInfo(h, h)
	if err != nil {
		return nil, 0, err
	}
	if bis.LastHeight < h {
		return nil, 0, errors.Errorf("Last height %d less than queries height %d", bis.LastHeight, h)
	}
	if len(bis.BlockMetas) != 1 {
		return nil, bis.LastHeight, errors.Errorf("Cannot get header for height %d", h)
	}
	// this is the header we actually want
	return bis.BlockMetas[0], bis.LastHeight, nil
}

// getCommit returns all the commit that proves the given
// block was approved by the validators.
//
// The current API requires we query block at h+1 to see the votes
// for block h
func (n Node) getCommit(h int) (*ttypes.Commit, error) {
	b, err := n.client.Block(h + 1)
	if err != nil {
		return nil, err
	}
	if b.Block == nil || b.Block.LastCommit == nil {
		return nil, errors.Errorf("No commit data for block %d", h+1)
	}
	commit := b.Block.LastCommit
	return commit, nil
}

func (n Node) validateCommitInfo(ci commitInfo) (lc.TmSignedHeader, error) {
	// now we validate it all and put it into our format for the
	res, err := validateHeaderInfo(ci.Header)
	if err != nil {
		return res, err
	}
	res.LastHeight = uint64(ci.MaxHeight)
	err = ci.Commit.ValidateBasic()
	if err != nil {
		return res, err
	}
	res.Votes, err = n.processVotes(ci.Commit.Precommits, res.Height(), res.Hash)
	return res, err
}

func validateHeaderInfo(header *ttypes.BlockMeta) (lc.TmSignedHeader, error) {
	var res lc.TmSignedHeader
	head := header.Header

	// make sure the hash matches
	calc := head.Hash()
	if !bytes.Equal(header.Hash, calc) {
		return res,
			errors.Errorf("Calculated header hash: %X, claimed header hash: %X",
				calc, header.Hash)
	}

	// this header looks good, transform the data!
	res = lc.TmSignedHeader{
		Hash: header.Hash,
		Header: lc.TmHeader{
			ChainID:        head.ChainID,
			Height:         uint64(head.Height),
			Time:           head.Time,
			LastBlockID:    head.LastBlockID.Hash,
			LastCommitHash: head.LastCommitHash,
			DataHash:       head.DataHash,
			ValidatorsHash: head.ValidatorsHash,
			AppHash:        head.AppHash,
		},
	}
	return res, nil
}

// processVotes is very similar to tendermint/types/validator_set.go:VerifyCommit
// and we should track changes there.  However, that code requires access to the
// cannonical validator set.  And here, we just want to get the data, so we can
// pass to a Certifier for processing. There will be various strategies for
// syncing the validator set in a Certifier, so we don't want to hard-code it here.
//
// also note that `err = b.Block.LastCommit.ValidateBasic()`
// in getCommit does a number of checks already, like they are all for the
// same block
func (n Node) processVotes(votes []*ttypes.Vote, height uint64, blockHash []byte) (lc.TmVotes, error) {
	res := make([]lc.TmVote, len(votes))
	h := int(height)

	// verify height and blockhash for the first vote (the rest are the same)
	if votes[0].Height != h {
		return nil, errors.New("Votes have incorrect height")
	}

	i := 0
	for _, v := range votes {
		// some votes may be nil, just skip them (tendermint/types/block.go:298)
		if v == nil {
			continue
		}
		if !bytes.Equal(blockHash, v.BlockID.Hash) {
			// Precommits has all votes.  but we can skip those for other blocks
			// FIXME: do we need to compare with the full BlockID???
			continue
		}
		sign := ttypes.SignBytes(n.chainID, v)
		// and store the info we care about
		res[i] = lc.TmVote{
			SignBytes:        sign,
			ValidatorAddress: v.ValidatorAddress,
			Signature:        v.Signature,
			Height:           uint64(v.Height),
			BlockHash:        v.BlockID.Hash,
		}

		// advance the count
		i++
	}

	return lc.TmVotes(res[:i]), nil
}
