package client

import (
	"fmt"

	"github.com/tendermint/go-wire/data"
	lc "github.com/tendermint/light-client"
	"github.com/tendermint/light-client/certifiers"
	"github.com/tendermint/light-client/proofs"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	"github.com/tendermint/tendermint/types"
	"github.com/tendermint/tmlibs/events"
)

var _ rpcclient.Client = Wrapper{}

type Wrapper struct {
	rpcclient.Client
	cert *certifiers.InquiringCertifier
}

func Wrap(c rpcclient.Client, cert *certifiers.InquiringCertifier) Wrapper {
	wrap := Wrapper{c, cert}
	// if we wrap http client, then we can swap out the event switch to filter
	if hc, ok := c.(*rpcclient.HTTP); ok {
		evt := hc.WSEvents.EventSwitch
		hc.WSEvents.EventSwitch = WrappedSwitch{evt, wrap}
	}
	return wrap
}

func (w Wrapper) ABCIQuery(path string, data data.Bytes, prove bool) (*ctypes.ResultABCIQuery, error) {
	r, err := w.Client.ABCIQuery(path, data, prove)
	if !prove || err != nil {
		return r, err
	}
	// get a verified commit to validate from
	c, err := w.Commit(int(r.Height))
	if err != nil {
		return nil, err
	}
	// make sure the checkpoint and proof match up
	check := lc.CheckpointFromResult(c)
	// verify query
	proof := proofs.AppProof{
		Height: r.Height,
		Key:    r.Key,
		Value:  r.Value,
		Proof:  r.Proof,
	}
	err = proof.Validate(check)
	return r, err
}

func (w Wrapper) Tx(hash []byte, prove bool) (*ctypes.ResultTx, error) {
	r, err := w.Client.Tx(hash, prove)
	if !prove || err != nil {
		return r, err
	}
	// get a verified commit to validate from
	c, err := w.Commit(r.Height)
	if err != nil {
		return nil, err
	}
	// make sure the checkpoint and proof match up
	check := lc.CheckpointFromResult(c)
	// verify tx
	proof := proofs.TxProof{
		Height: uint64(r.Height),
		Proof:  r.Proof,
	}
	err = proof.Validate(check)
	return r, err
}

func (w Wrapper) BlockchainInfo(minHeight, maxHeight int) (*ctypes.ResultBlockchainInfo, error) {
	r, err := w.Client.BlockchainInfo(minHeight, maxHeight)
	if err != nil {
		return nil, err
	}

	// go and verify every blockmeta in the result....
	for _, meta := range r.BlockMetas {
		// get a checkpoint to verify from
		c, err := w.Commit(meta.Header.Height)
		if err != nil {
			return nil, err
		}
		check := lc.CheckpointFromResult(c)
		err = proofs.ValidateBlockMeta(meta, check)
		if err != nil {
			return nil, err
		}
	}

	return r, nil
}

func (w Wrapper) Block(height int) (*ctypes.ResultBlock, error) {
	r, err := w.Client.Block(height)
	if err != nil {
		return nil, err
	}
	// get a checkpoint to verify from
	c, err := w.Commit(height)
	if err != nil {
		return nil, err
	}
	check := lc.CheckpointFromResult(c)

	// now verify
	err = proofs.ValidateBlockMeta(r.BlockMeta, check)
	if err != nil {
		return nil, err
	}
	err = proofs.ValidateBlock(r.Block, check)
	if err != nil {
		return nil, err
	}
	return r, nil
}

// Commit downloads the Commit and certifies it with the certifiers.
//
// This is the foundation for all other verification in this module
func (w Wrapper) Commit(height int) (*ctypes.ResultCommit, error) {
	rpcclient.WaitForHeight(w.Client, height, nil)
	r, err := w.Client.Commit(height)
	// if we got it, then certify it
	if err == nil {
		check := lc.CheckpointFromResult(r)
		err = w.cert.Certify(check)
	}
	return r, err
}

type WrappedSwitch struct {
	types.EventSwitch
	client rpcclient.Client
}

func (s WrappedSwitch) FireEvent(event string, data events.EventData) {
	tm, ok := data.(types.TMEventData)
	if !ok {
		fmt.Printf("bad type %#v\n", data)
		return
	}

	// check to validate it if possible, and drop if not valid
	switch t := tm.Unwrap().(type) {
	case types.EventDataNewBlockHeader:
		err := verifyHeader(s.client, t.Header)
		if err != nil {
			fmt.Printf("Invalid header: %#v\n", err)
			return
		}
	case types.EventDataNewBlock:
		err := verifyBlock(s.client, t.Block)
		if err != nil {
			fmt.Printf("Invalid block: %#v\n", err)
			return
		}
	}

	// looks good, we fire it
	s.EventSwitch.FireEvent(event, data)
}

func verifyHeader(c rpcclient.Client, head *types.Header) error {
	// get a checkpoint to verify from
	commit, err := c.Commit(head.Height)
	if err != nil {
		return err
	}
	check := lc.CheckpointFromResult(commit)
	return proofs.ValidateHeader(head, check)
}

func verifyBlock(c rpcclient.Client, block *types.Block) error {
	// get a checkpoint to verify from
	commit, err := c.Commit(block.Height)
	if err != nil {
		return err
	}
	check := lc.CheckpointFromResult(commit)
	return proofs.ValidateBlock(block, check)
}
