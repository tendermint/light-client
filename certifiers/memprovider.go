package certifiers

import (
	"encoding/hex"
	"sort"

	certerr "github.com/tendermint/light-client/certifiers/errors"
)

type memStoreProvider struct {
	// byHeight is always sorted by Height... need to support range search (nil, h]
	// btree would be more efficient for larger sets
	byHeight seeds
	byHash   map[string]Seed
}

func NewMemStoreProvider() Provider {
	return &memStoreProvider{
		byHeight: seeds{},
		byHash:   map[string]Seed{},
	}
}

func (m *memStoreProvider) encodeHash(hash []byte) string {
	return hex.EncodeToString(hash)
}

func (m *memStoreProvider) StoreSeed(seed Seed) error {
	// make sure the seed is self-consistent before saving
	err := seed.ValidateBasic(seed.Commit.Header.ChainID)
	if err != nil {
		return err
	}

	// store the valid seed
	key := m.encodeHash(seed.Hash())
	m.byHash[key] = seed
	m.byHeight = append(m.byHeight, seed)
	sort.Sort(m.byHeight)
	return nil
}

func (m *memStoreProvider) GetByHeight(h int) (Seed, error) {
	// search from highest to lowest
	for i := len(m.byHeight) - 1; i >= 0; i-- {
		s := m.byHeight[i]
		if s.Height() <= h {
			return s, nil
		}
	}
	return Seed{}, certerr.ErrSeedNotFound()
}

func (m *memStoreProvider) GetByHash(hash []byte) (Seed, error) {
	var err error
	s, ok := m.byHash[m.encodeHash(hash)]
	if !ok {
		err = certerr.ErrSeedNotFound()
	}
	return s, err
}

func (m *memStoreProvider) LatestSeed() (Seed, error) {
	l := len(m.byHeight)
	if l == 0 {
		return Seed{}, certerr.ErrSeedNotFound()
	}
	return m.byHeight[l-1], nil
}
