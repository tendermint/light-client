package certifiers

import (
	lc "github.com/tendermint/light-client"
	"github.com/tendermint/tendermint/types"
)

type Seed struct {
	lc.Checkpoint
	Validators []*types.Validator
}

// Provider is used to get more validators by other means
type Provider interface {
	StoreSeed(Seed) error
	GetByHeight(int) (Seed, error)
	GetByHash([]byte) (Seed, error)
}

type CacheProvider struct {
	Cache  Provider
	Source Provider
}

func (c CacheProvider) StoreSeed(seed Seed) error {
	c.Source.StoreSeed(seed)
	return c.Cache.StoreSeed(seed)
}

func (c CacheProvider) GetByHeight(h int) (Seed, error) {
	s, err := c.Cache.GetByHeight(h)
	if err != nil {
		s, err = c.Source.GetByHeight(h)
		if err == nil {
			c.Cache.StoreSeed(s)
		}
	}
	return s, err
}

func (c CacheProvider) GetByHash(hash []byte) (Seed, error) {
	s, err := c.Cache.GetByHash(hash)
	if err != nil {
		s, err = c.Source.GetByHash(hash)
		if err == nil {
			c.Cache.StoreSeed(s)
		}
	}
	return s, err
}

// Also MemStoreProvider, FileStoreProvider, NodeProvider, ...
