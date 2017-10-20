package files

import (
	"encoding/json"
	"os"

	"github.com/pkg/errors"

	wire "github.com/tendermint/go-wire"

	"github.com/tendermint/light-client/certifiers"
	certerr "github.com/tendermint/light-client/certifiers/errors"
)

const (
	MaxSeedSize = 1024 * 1024
)

// SaveSeed exports the seed in binary / go-wire style
func SaveSeed(s certifiers.Seed, path string) error {
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

// SaveSeedJSON exports the seed in a json format
func SaveSeedJSON(s certifiers.Seed, path string) error {
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

func LoadSeed(path string) (certifiers.Seed, error) {
	var seed certifiers.Seed
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

func LoadSeedJSON(path string) (certifiers.Seed, error) {
	var seed certifiers.Seed
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
