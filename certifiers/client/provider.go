package client

import (
	"bytes"

	rpcclient "github.com/tendermint/tendermint/rpc/client"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	"github.com/tendermint/tendermint/types"

	"github.com/tendermint/light-client/certifiers"
	certerr "github.com/tendermint/light-client/certifiers/errors"
)

var _ certifiers.Provider = &Provider{}

type SignStatusClient interface {
	rpcclient.SignClient
	rpcclient.StatusClient
}

type Provider struct {
	node       SignStatusClient
	lastHeight int
}

func New(node SignStatusClient) *Provider {
	return &Provider{node: node}
}

func NewHTTP(remote string) *Provider {
	return &Provider{
		node: rpcclient.NewHTTP(remote, "/websocket"),
	}
}

// StoreSeed is a noop, as clients can only read from the chain...
func (p *Provider) StoreSeed(_ certifiers.Seed) error { return nil }

// GetHash gets the most recent validator and sees if it matches
//
// TODO: improve when the rpc interface supports more functionality
func (p *Provider) GetByHash(hash []byte) (certifiers.Seed, error) {
	var seed certifiers.Seed
	vals, err := p.node.Validators(nil)
	// if we get no validators, or a different height, return an error
	if err != nil {
		return seed, err
	}
	p.updateHeight(vals.BlockHeight)
	vhash := types.NewValidatorSet(vals.Validators).Hash()
	if !bytes.Equal(hash, vhash) {
		return seed, certerr.ErrSeedNotFound()
	}
	return p.seedFromVals(vals)
}

// GetByHeight gets the validator set by height
func (p *Provider) GetByHeight(h int) (seed certifiers.Seed, err error) {
	commit, err := p.node.Commit(&h)
	if err != nil {
		return seed, err
	}
	return p.seedFromCommit(commit)
}

func (p *Provider) LatestSeed() (seed certifiers.Seed, err error) {
	commit, err := p.GetLatestCommit()
	if err != nil {
		return seed, err
	}
	return p.seedFromCommit(commit)
}

// GetLatestCommit should return the most recent commit there is,
// which handles queries for future heights as per the semantics
// of GetByHeight.
func (p *Provider) GetLatestCommit() (*ctypes.ResultCommit, error) {
	status, err := p.node.Status()
	if err != nil {
		return nil, err
	}
	return p.node.Commit(&status.LatestBlockHeight)
}

func (p *Provider) seedFromVals(vals *ctypes.ResultValidators) (certifiers.Seed, error) {
	seed := certifiers.Seed{
		Validators: types.NewValidatorSet(vals.Validators),
	}
	// now get the commits and build a seed
	commit, err := p.node.Commit(&vals.BlockHeight)
	if err != nil {
		return seed, err
	}
	seed.Commit = certifiers.CommitFromResult(commit)
	return seed, nil
}

func (p *Provider) seedFromCommit(commit *ctypes.ResultCommit) (certifiers.Seed, error) {
	seed := certifiers.Seed{
		Commit: certifiers.CommitFromResult(commit),
	}

	// now get the proper validators
	vals, err := p.node.Validators(&commit.Header.Height)
	if err != nil {
		return seed, err
	}

	// make sure they match the commit (as we cannot enforce height)
	vset := types.NewValidatorSet(vals.Validators)
	if !bytes.Equal(vset.Hash(), commit.Header.ValidatorsHash) {
		return seed, certerr.ErrValidatorsChanged()
	}

	p.updateHeight(commit.Header.Height)
	seed.Validators = vset
	return seed, nil
}

func (p *Provider) updateHeight(h int) {
	if h > p.lastHeight {
		p.lastHeight = h
	}
}
