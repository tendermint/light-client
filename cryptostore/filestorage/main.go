/*
package filestorage provides a secure on-disk storage of private keys and
metadata.  Security is enforced by file and directory permissions, much
like standard ssh key storage.
*/
package filestorage

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/pkg/errors"
	crypto "github.com/tendermint/go-crypto"
	lightclient "github.com/tendermint/light-client"
	"github.com/tendermint/light-client/cryptostore"
)

const (
	BlockType = "Tendermint Light Client"
	privExt   = "tlc"
	pubExt    = "pub"
	keyPerm   = os.FileMode(0600)
	pubPerm   = os.FileMode(0644)
	dirPerm   = os.FileMode(0700)
)

type store struct {
	keyDir string
}

// New creates an instance of file-based key storage with tight permissions
func New(dir string) cryptostore.Storage {
	err := os.Mkdir(dir, dirPerm)
	if err != nil && !os.IsExist(err) {
		panic(err)
	}
	return store{dir}
}

// Put creates two files, one with the public info as json, the other
// with the (encoded) private key as gpg ascii-armor style
func (s store) Put(name string, key []byte, info lightclient.KeyInfo) error {
	pub, priv := s.nameToPaths(name)

	// write public info
	err := writeInfo(pub, info)
	if err != nil {
		return err
	}

	// write private info
	return writeKey(priv, name, key)
}

// Get loads the keyinfo and (encoded) private key from the directory
// It uses `name` to generate the filename, and returns an error if the
// files don't exist or are in the incorrect format
func (s store) Get(name string) ([]byte, lightclient.KeyInfo, error) {
	pub, priv := s.nameToPaths(name)

	info, err := readInfo(pub)
	if err != nil {
		return nil, info, err
	}

	key, err := readKey(priv)
	return key, info, err
}

// List parses the key directory for public info and returns a list of
// KeyInfo for all keys located in this directory.
func (s store) List() ([]lightclient.KeyInfo, error) {
	dir, err := os.Open(s.keyDir)
	if err != nil {
		return nil, errors.Wrap(err, "List Keys")
	}
	names, err := dir.Readdirnames(0)
	if err != nil {
		return nil, errors.Wrap(err, "List Keys")
	}

	// filter names for .pub ending and load them one by one
	// half the files is a good guess for pre-allocating the slice
	infos := make([]lightclient.KeyInfo, len(names)/2, 0)
	for _, name := range names {
		p := path.Join(s.keyDir, name)
		info, err := readInfo(p)
		if err != nil {
			return nil, err
		}
		infos = append(infos, info)
	}

	return infos, nil
}

// Delete permanently removes the public and private info for the named key
// The calling function should provide some security checks first.
func (s store) Delete(name string) error {
	pub, priv := s.nameToPaths(name)
	err := os.Remove(priv)
	if err != nil {
		return errors.Wrap(err, "Deleting Private Key")
	}
	err = os.Remove(pub)
	return errors.Wrap(err, "Deleting Public Key")
}

func (s store) nameToPaths(name string) (pub, priv string) {
	privName := fmt.Sprintf("%s.%s", name, privExt)
	pubName := fmt.Sprintf("%s.%s", name, pubExt)
	return path.Join(s.keyDir, privName), path.Join(s.keyDir, pubName)
}

func readInfo(path string) (lightclient.KeyInfo, error) {
	var info lightclient.KeyInfo
	f, err := os.Open(path)
	if err != nil {
		return info, errors.Wrap(err, "Getting Public Key")
	}
	data, err := ioutil.ReadAll(f)
	if err != nil {
		return info, errors.Wrap(err, "Getting Public Key")
	}
	err = json.Unmarshal(data, &info)
	return info, errors.Wrap(err, "Parsing Public Key")
}

func readKey(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrap(err, "Getting Private Key")
	}
	d, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, errors.Wrap(err, "Getting Private Key")
	}
	block, _, key, err := crypto.DecodeArmor(string(d))
	if err != nil {
		return nil, errors.Wrap(err, "Invalid Armor")
	}
	if block != BlockType {
		return nil, errors.Errorf("Unknown key type: %s", block)
	}
	return key, nil
}

func writeInfo(path string, info lightclient.KeyInfo) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, pubPerm)
	if err != nil {
		return errors.Wrap(err, "Saving Public Key")
	}
	defer f.Close()
	d, err := json.Marshal(info)
	if err != nil {
		return errors.Wrap(err, "Saving Public Key")
	}
	_, err = f.Write(d)
	return errors.Wrap(err, "Saving Public Key")
}

func writeKey(path, name string, key []byte) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, keyPerm)
	if err != nil {
		return errors.Wrap(err, "Saving Private Key")
	}
	defer f.Close()
	headers := map[string]string{"name": name}
	text := crypto.EncodeArmor(BlockType, headers, key)
	_, err = f.WriteString(text)
	return errors.Wrap(err, "Saving Private Key")
}
