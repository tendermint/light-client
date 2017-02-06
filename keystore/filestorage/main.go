package filestorage

import (
	"fmt"
	"os"
	"path"

	lightclient "github.com/tendermint/light-client"
	"github.com/tendermint/light-client/keystore"
)

const (
	BlockType = "Tendermint Light Client"
	keyPerm   = os.FileMode(0600)
	dirPerm   = os.FileMode(0700)
)

type store struct {
	keyDir string
}

// TODO: implement
func New(dir string) keystore.Storage {
	err := os.Mkdir(dir, dirPerm)
	if err != nil && !os.IsExist(err) {
		panic(err)
	}
	return store{dir}
}

func (s store) Put(name string, key []byte, info lightclient.KeyInfo) error {
	return nil
}

func (s store) Get(name string) ([]byte, lightclient.KeyInfo, error) {
	return nil, lightclient.KeyInfo{}, nil
}

func (s store) List() ([]lightclient.KeyInfo, error) {
	return nil, nil
}

func (s store) Delete(name string) error {
	return nil
}

func (s store) nameToPath(name string) string {
	fname := fmt.Sprintf("%s.tlc", name)
	return path.Join(s.keyDir, fname)
}
