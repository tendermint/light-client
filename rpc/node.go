package rpc

import (
	lc "github.com/tendermint/light-client"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
)

type Node struct {
	client *HTTPClient
}

func NewNode(rpcAddr string) Node {
	return Node{
		client: NewClient(rpcAddr, "/websocket"),
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
	return queryResp(qr), err
}

// Prove returns a merkle proof for the given key
func (n Node) Prove(key []byte) (lc.TmQueryResult, error) {
	qr, err := n.client.ABCIQuery("/key", key, true)
	return queryResp(qr), err
}

func queryResp(qr *ctypes.ResultABCIQuery) lc.TmQueryResult {
	if qr == nil {
		return lc.TmQueryResult{}
	}
	r := qr.Response
	return lc.TmQueryResult{
		Code:   lc.TmCode(r.Code),
		Key:    r.Key,
		Value:  r.Value,
		Proof:  r.Proof,
		Log:    r.Log,
		Height: r.Height,
	}
}

// SignedHeader gives us Header data along with the backing signatures
//
// It is also responsible for making the blockhash and header match,
// and that all votes are valid pre-commit votes for this block
// It does not check if the keys signing the votes are actual validators
//
// TODO
func (n Node) SignedHeader(height uint64) (lc.TmSignedHeader, error) {
	return lc.TmSignedHeader{}, nil
}
