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
	MaxFullCommitSize = 1024 * 1024
)

// SaveFullCommit exports the seed in binary / go-wire style
func SaveFullCommit(fc certifiers.FullCommit, path string) error {
	f, err := os.Create(path)
	if err != nil {
		// if os.IsExist(err) {
		//   return nil
		// }
		return errors.WithStack(err)
	}
	defer f.Close()

	var n int
	wire.WriteBinary(fc, f, &n, &err)
	return errors.WithStack(err)
}

// SaveFullCommitJSON exports the seed in a json format
func SaveFullCommitJSON(fc certifiers.FullCommit, path string) error {
	f, err := os.Create(path)
	if err != nil {
		// if os.IsExist(err) {
		//   return nil
		// }
		return errors.WithStack(err)
	}
	defer f.Close()
	stream := json.NewEncoder(f)
	err = stream.Encode(fc)
	return errors.WithStack(err)
}

func LoadFullCommit(path string) (certifiers.FullCommit, error) {
	var fc certifiers.FullCommit
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fc, certerr.ErrFullCommitNotFound()
		}
		return fc, errors.WithStack(err)
	}
	defer f.Close()

	var n int
	wire.ReadBinaryPtr(&fc, f, MaxFullCommitSize, &n, &err)
	return fc, errors.WithStack(err)
}

func LoadFullCommitJSON(path string) (certifiers.FullCommit, error) {
	var fc certifiers.FullCommit
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fc, certerr.ErrFullCommitNotFound()
		}
		return fc, errors.WithStack(err)
	}
	defer f.Close()

	stream := json.NewDecoder(f)
	err = stream.Decode(&fc)
	return fc, errors.WithStack(err)
}
