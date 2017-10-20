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
	byHash   map[string]FullCommit
}

// seeds just exists to allow easy sorting
type seeds []FullCommit

func (s seeds) Len() int      { return len(s) }
func (s seeds) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s seeds) Less(i, j int) bool {
	return s[i].Height() < s[j].Height()
}

func NewMemStoreProvider() Provider {
	return &memStoreProvider{
		byHeight: seeds{},
		byHash:   map[string]FullCommit{},
	}
}

func (m *memStoreProvider) encodeHash(hash []byte) string {
	return hex.EncodeToString(hash)
}

func (m *memStoreProvider) StoreFullCommit(seed FullCommit) error {
	// make sure the seed is self-consistent before saving
	err := seed.ValidateBasic(seed.Commit.Header.ChainID)
	if err != nil {
		return err
	}

	// store the valid seed
	key := m.encodeHash(seed.ValidatorsHash())
	m.byHash[key] = seed
	m.byHeight = append(m.byHeight, seed)
	sort.Sort(m.byHeight)
	return nil
}

func (m *memStoreProvider) GetByHeight(h int) (FullCommit, error) {
	// search from highest to lowest
	for i := len(m.byHeight) - 1; i >= 0; i-- {
		s := m.byHeight[i]
		if s.Height() <= h {
			return s, nil
		}
	}
	return FullCommit{}, certerr.ErrFullCommitNotFound()
}

func (m *memStoreProvider) GetByHash(hash []byte) (FullCommit, error) {
	var err error
	s, ok := m.byHash[m.encodeHash(hash)]
	if !ok {
		err = certerr.ErrFullCommitNotFound()
	}
	return s, err
}

func (m *memStoreProvider) LatestFullCommit() (FullCommit, error) {
	l := len(m.byHeight)
	if l == 0 {
		return FullCommit{}, certerr.ErrFullCommitNotFound()
	}
	return m.byHeight[l-1], nil
}
