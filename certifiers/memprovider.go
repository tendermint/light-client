package certifiers

import (
	"encoding/hex"
	"sort"
)

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
	// make sure the seed is self-consistent before saving
	err := seed.ValidateBasic(seed.Checkpoint.Header.ChainID)
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

func (m *MemStoreProvider) GetByHeight(h int) (Seed, error) {
	// search from highest to lowest
	for i := len(m.byHeight) - 1; i >= 0; i-- {
		s := m.byHeight[i]
		if s.Height() <= h {
			return s, nil
		}
	}
	return Seed{}, ErrSeedNotFound()
}

func (m *MemStoreProvider) GetByHash(hash []byte) (Seed, error) {
	var err error
	s, ok := m.byHash[m.encodeHash(hash)]
	if !ok {
		err = ErrSeedNotFound()
	}
	return s, err
}
