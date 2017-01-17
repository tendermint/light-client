package keystore

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/pkg/errors"
	crypto "github.com/tendermint/go-crypto"
	wire "github.com/tendermint/go-wire"
)

const (
	BlockType = "Tendermint Light Client"
	keyPerm   = os.FileMode(0600)
	dirPerm   = os.FileMode(0700)
)

type Store struct {
	keyDir string
}

// TODO: error not panic?
func New(dir string) Store {
	// make sure the dir is there...
	err := os.Mkdir(dir, dirPerm)
	if err != nil && !os.IsExist(err) {
		panic(err)
	}
	return Store{keyDir: dir}
}

func secret(passphrase string) []byte {
	// TODO: Sha256(Bcrypt(passphrase))
	return crypto.Sha256([]byte(passphrase))
}

func generateKey(name, passphrase string) string {
	priv := crypto.GenPrivKeyEd25519()
	secret := secret(passphrase)
	// TODO: crypto.MixEntropy()
	cipher := crypto.EncryptSymmetric(priv.Bytes(), secret)
	headers := map[string]string{"name": name}
	text := crypto.EncodeArmor(BlockType, headers, cipher)
	return text
}

func decryptKey(cipher, passphrase string) (crypto.PrivKey, error) {
	block, _, data, err := crypto.DecodeArmor(cipher)
	if err != nil {
		return nil, errors.Wrap(err, "Invalid Armor")
	}
	if block != BlockType {
		return nil, errors.Errorf("Unknown key type: %s", block)
	}

	secret := secret(passphrase)
	private, err := crypto.DecryptSymmetric(data, secret)
	if err != nil {
		return nil, errors.Wrap(err, "Invalid Passphrase")
	}

	key, err := crypto.PrivKeyFromBytes(private)
	if err != nil {
		return nil, errors.Wrap(err, "Invalid Passphrase")
	}
	return key, nil
}

func (s Store) nameToPath(name string) string {
	fname := fmt.Sprintf("%s.tlc", name)
	return path.Join(s.keyDir, fname)
}

func (s Store) CreateKey(name, passphrase string) error {
	path := s.nameToPath(name)
	file, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, keyPerm)
	if err != nil {
		return err
	}
	defer file.Close()

	text := generateKey(name, passphrase)
	_, err = file.WriteString(text)
	return err
}

func (s Store) SignMessage(msg []byte, name, passphrase string) (sig []byte, pubkey []byte, err error) {
	path := s.nameToPath(name)
	file, err := os.Open(path)
	if err != nil {
		return nil, nil, errors.Wrap(err, "Cannot open keyfile")
	}

	keyData, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, nil, errors.Wrap(err, "Cannot read keyfile")
	}

	key, err := decryptKey(string(keyData), passphrase)
	if err != nil {
		return nil, nil, err
	}

	// go-wire format to be compatible with basecoin, and likely others
	sig = wire.BinaryBytes(key.Sign(msg))
	pubkey = wire.BinaryBytes(key.PubKey())
	return sig, pubkey, nil
}
