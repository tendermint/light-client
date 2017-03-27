package certifiers

import (
	lc "github.com/tendermint/light-client"
	"github.com/tendermint/tendermint/types"
)

// Seed is a checkpoint and the actual validator set, the base info you
// need to update to a given point, assuming knowledge of some previous
// validator set
type Seed struct {
	lc.Checkpoint
	Validators []*types.Validator
}

func (s Seed) Height() int {
	return s.Checkpoint.Height()
}

func (s Seed) Hash() []byte {
	return s.Checkpoint.Header.ValidatorsHash
}

type Seeds []Seed

func (s Seeds) Len() int      { return len(s) }
func (s Seeds) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s Seeds) Less(i, j int) bool {
	return s[i].Height() < s[j].Height()
}

// Provider is used to get more validators by other means
//
// TODO: Also FileStoreProvider, NodeProvider, ...
type Provider interface {
	StoreSeed(seed Seed) error
	// GetByHeight returns the closest seed at with height <= h
	GetByHeight(h int) (Seed, error)
	// GetByHash returns a seed exactly matching this validator hash
	GetByHash(hash []byte) (Seed, error)
}

// CacheProvider allows you to place one or more caches in front of a source
// Provider.  It runs through them in order until a match is found.
// So you can keep a local cache, and check with the network if
// no data is there.
type CacheProvider struct {
	Providers []Provider
}

func NewCacheProvider(providers ...Provider) CacheProvider {
	return CacheProvider{
		Providers: providers,
	}
}

// StoreSeed tries to add the seed to all providers.
//
// Aborts on first error it encounters (closest provider)
func (c CacheProvider) StoreSeed(seed Seed) (err error) {
	for _, p := range c.Providers {
		err := p.StoreSeed(seed)
		if err != nil {
			break
		}
	}
	return err
}

func (c CacheProvider) GetByHeight(h int) (s Seed, err error) {
	for _, p := range c.Providers {
		s, err = p.GetByHeight(h)
		if err == nil {
			break
		}
	}
	return s, err
}

func (c CacheProvider) GetByHash(hash []byte) (s Seed, err error) {
	for _, p := range c.Providers {
		s, err = p.GetByHash(hash)
		if err == nil {
			break
		}
	}
	return s, err
}

// MissingProvider doens't store anything, always a miss
// Designed as a mock for testing
type MissingProvider struct{}

func NewMissingProvider() MissingProvider {
	return MissingProvider{}
}

func (_ MissingProvider) StoreSeed(_ Seed) error { return nil }
func (_ MissingProvider) GetByHeight(_ int) (Seed, error) {
	return Seed{}, ErrSeedNotFound()
}
func (_ MissingProvider) GetByHash(_ []byte) (Seed, error) {
	return Seed{}, ErrSeedNotFound()
}
