package certifiers

import (
	"encoding/json"
	"os"

	"github.com/pkg/errors"

	wire "github.com/tendermint/go-wire"

	"github.com/tendermint/tendermint/types"

	certerr "github.com/tendermint/light-client/certifiers/errors"
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

// Write exports the seed in binary / go-wire style
func (s Seed) Write(path string) error {
	f, err := os.Create(path)
	if err != nil {
		// if os.IsExist(err) {
		//   return nil
		// }
		return errors.WithStack(err)
	}
	defer f.Close()

	var n int
	wire.WriteBinary(s, f, &n, &err)
	return errors.WithStack(err)
}

// Write exports the seed in a json format
func (s Seed) WriteJSON(path string) error {
	f, err := os.Create(path)
	if err != nil {
		// if os.IsExist(err) {
		//   return nil
		// }
		return errors.WithStack(err)
	}
	defer f.Close()
	stream := json.NewEncoder(f)
	err = stream.Encode(s)
	return errors.WithStack(err)
}

func LoadSeed(path string) (Seed, error) {
	var seed Seed
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return seed, certerr.ErrSeedNotFound()
		}
		return seed, errors.WithStack(err)
	}
	defer f.Close()

	var n int
	wire.ReadBinaryPtr(&seed, f, MaxSeedSize, &n, &err)
	return seed, errors.WithStack(err)
}

func LoadSeedJSON(path string) (Seed, error) {
	var seed Seed
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return seed, certerr.ErrSeedNotFound()
		}
		return seed, errors.WithStack(err)
	}
	defer f.Close()

	stream := json.NewDecoder(f)
	err = stream.Decode(&seed)
	return seed, errors.WithStack(err)
}

// seeds just exists to allow easy sorting
type seeds []Seed

func (s seeds) Len() int      { return len(s) }
func (s seeds) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s seeds) Less(i, j int) bool {
	return s[i].Height() < s[j].Height()
}
