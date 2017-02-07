/*
package memstorage provides a simple in-memory key store designed for
use in test cases, particularly to isolate them from the filesystem and
concurrency, cleanup issues.
*/
package memstorage

import (
	"github.com/pkg/errors"
	lightclient "github.com/tendermint/light-client"
	"github.com/tendermint/light-client/cryptostore"
)

type data struct {
	info lightclient.KeyInfo
	key  []byte
}

type store map[string]data

// New creates an instance of file-based key storage with tight permissions
func New(dir string) cryptostore.Storage {
	return store{}
}

func (s store) Put(name string, key []byte, info lightclient.KeyInfo) error {
	if _, ok := s[name]; ok {
		return errors.Errorf("Key named '%s' already exists", name)
	}
	s[name] = data{info, key}
	return nil
}

func (s store) Get(name string) ([]byte, lightclient.KeyInfo, error) {
	var err error
	d, ok := s[name]
	if !ok {
		err = errors.Errorf("Key named '%s' doesn't exist", name)
	}
	return d.key, d.info, err
}

func (s store) List() ([]lightclient.KeyInfo, error) {
	res := make([]lightclient.KeyInfo, len(s))
	i := 0
	for _, d := range s {
		res[i] = d.info
		i++
	}
	return res, nil
}

func (s store) Delete(name string) error {
	_, ok := s[name]
	if !ok {
		return errors.Errorf("Key named '%s' doesn't exist", name)
	}
	delete(s, name)
	return nil
}
