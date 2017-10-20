package certifiers

import (
	"github.com/tendermint/tendermint/types"
)

const (
	MaxSeedSize = 1024 * 1024
)

// Seed is a commit and the actual validator set, the base info you
// need to update to a given point, assuming knowledge of some previous
// validator set
type Seed struct {
	*Commit    `json:"commit"`
	Validators *types.ValidatorSet `json:"validator_set"`
}

func (s Seed) Height() int {
	if s.Commit == nil {
		return 0
	}
	return s.Commit.Height()
}

func (s Seed) Hash() []byte {
	h := s.Commit.Header
	if h == nil {
		return nil
	}
	return h.ValidatorsHash
}

// seeds just exists to allow easy sorting
type seeds []Seed

func (s seeds) Len() int      { return len(s) }
func (s seeds) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s seeds) Less(i, j int) bool {
	return s[i].Height() < s[j].Height()
}
