package client

import (
	"bytes"

	"github.com/pkg/errors"
	"github.com/tendermint/light-client/certifiers"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	"github.com/tendermint/tendermint/types"
)

var _ certifiers.Provider = &Provider{}

type Provider struct {
	node       rpcclient.SignClient
	lastHeight int
}

func New(node rpcclient.SignClient) *Provider {
	return &Provider{node: node}
}

func NewHTTP(remote string) *Provider {
	return &Provider{
		node: rpcclient.NewHTTP(remote, "/websocket"),
	}
}

// StoreSeed is a noop, as clients can only read from the chain...
func (p *Provider) StoreSeed(_ certifiers.Seed) error { return nil }

// GetHash gets the most recent validator (only one available)
// and sees if it matches
//
// TODO: improve when the rpc interface supports more functionality
func (p *Provider) GetByHash(hash []byte) (certifiers.Seed, error) {
	var seed certifiers.Seed
	vals, err := p.node.Validators()
	// if we get no validators, or a different height, return an error
	if err != nil {
		return seed, errors.WithStack(err)
	}
	p.updateHeight(vals.BlockHeight)
	vhash := types.NewValidatorSet(vals.Validators).Hash()
	if !bytes.Equal(hash, vhash) {
		return seed, certifiers.ErrSeedNotFound()
	}
	return p.buildSeed(vals)
}

// GetByHeight gets the most recent validator (only one available)
// and sees if it matches
//
// TODO: keep track of most recent height, it will never go down
//
// TODO: improve when the rpc interface supports more functionality
func (p *Provider) GetByHeight(h int) (certifiers.Seed, error) {
	var seed certifiers.Seed
	vals, err := p.node.Validators()
	// if we get no validators, or a different height, return an error
	if err != nil {
		return seed, errors.WithStack(err)
	}
	p.updateHeight(vals.BlockHeight)
	if vals.BlockHeight > h {
		return seed, certifiers.ErrSeedNotFound()
	}
	return p.buildSeed(vals)
}

func (p *Provider) buildSeed(vals *ctypes.ResultValidators) (certifiers.Seed, error) {
	seed := certifiers.Seed{
		Validators: types.NewValidatorSet(vals.Validators),
	}
	// looks good, now get the commits and build a seed
	commit, err := p.node.Commit(vals.BlockHeight)
	if err == nil {
		seed.Header = commit.Header
		seed.Commit = commit.Commit
	}
	return seed, errors.WithStack(err)
}

func (p *Provider) updateHeight(h int) {
	if h > p.lastHeight {
		p.lastHeight = h
	}
}
