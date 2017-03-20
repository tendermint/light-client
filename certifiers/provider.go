package certifiers

import (
	"encoding/hex"
	rawerr "errors"
	"sort"

	"github.com/pkg/errors"
	lc "github.com/tendermint/light-client"
	"github.com/tendermint/tendermint/types"
)

var (
	errSeedNotFound = rawerr.New("Seed not found by provider")
)

// SeedNotFound asserts whether an error is due to missing data
func SeedNotFound(err error) bool {
	return err != nil && (errors.Cause(err) == errSeedNotFound)
}

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

type MemStoreProvider struct {
	// byHeight is always sorted by Height... need to support range search (nil, h]
	// btree would be more efficient for larger sets
	byHeight Seeds
	byHash   map[string]Seed
}

func NewMemStoreProvider() *MemStoreProvider {
	return &MemStoreProvider{
		byHeight: Seeds{},
		byHash:   map[string]Seed{},
	}
}

func (m *MemStoreProvider) encodeHash(hash []byte) string {
	return hex.EncodeToString(hash)
}

func (m *MemStoreProvider) StoreSeed(seed Seed) error {
	key := m.encodeHash(seed.Hash())
	m.byHash[key] = seed
	m.byHeight = append(m.byHeight, seed)
	sort.Sort(m.byHeight)
	return nil
}

func (m *MemStoreProvider) GetByHeight(h int) (Seed, error) {
	// search from highest to lowest
	for i := len(m.byHeight) - 1; i >= 0; i-- {
		s := m.byHeight[i]
		if s.Height() <= h {
			return s, nil
		}
	}
	return Seed{}, errors.WithStack(errSeedNotFound)
}

func (m *MemStoreProvider) GetByHash(hash []byte) (Seed, error) {
	var err error
	s, ok := m.byHash[m.encodeHash(hash)]
	if !ok {
		err = errors.WithStack(errSeedNotFound)
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
	return Seed{}, errors.WithStack(errSeedNotFound)
}
func (_ MissingProvider) GetByHash(_ []byte) (Seed, error) {
	return Seed{}, errors.WithStack(errSeedNotFound)
}
