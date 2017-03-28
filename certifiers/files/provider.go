package files

import (
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/pkg/errors"
	wire "github.com/tendermint/go-wire"
	"github.com/tendermint/light-client/certifiers"
)

const (
	Ext      = ".tsd"
	ValDir   = "validators"
	CheckDir = "checkpoints"
	dirPerm  = os.FileMode(0755)
	filePerm = os.FileMode(0644)
)

// Provider stores all data in the filesystem
// We assume the same validator hash may be reused by many different
// headers/Checkpoints, and thus store it separately. This leaves us
// with three issues:
//
// 1. Given a validator hash, retrieve the validator set if previously stored
// 2. Given a block height, find the Checkpoint with the highest height <= h
// 3. Given a Seed, store it quickly to satisfy 1 and 2
//
// Note that we do not worry about caching, as that can be achieved by
// pairing this with a MemStoreProvider and CacheProvider from certifiers
type Provider struct {
	valDir   string
	checkDir string
}

// NewProvider creates the parent dir and subdirs
// for validators and checkpoints as needed
func NewProvider(dir string) Provider {
	valDir := filepath.Join(dir, ValDir)
	checkDir := filepath.Join(dir, CheckDir)
	for _, d := range []string{valDir, checkDir} {
		err := os.MkdirAll(d, dirPerm)
		if err != nil {
			panic(err)
		}
	}
	return Provider{valDir: valDir, checkDir: checkDir}
}

func (m Provider) encodeHash(hash []byte) string {
	return hex.EncodeToString(hash) + Ext
}

func (m Provider) encodeHeight(h int) string {
	// pad up to 10^12 for height...
	return fmt.Sprintf("%012d%s", h, Ext)
}

func (m Provider) StoreSeed(seed certifiers.Seed) error {
	// make sure the seed is self-consistent before saving
	err := seed.ValidateBasic(seed.Checkpoint.Header.ChainID)
	if err != nil {
		return err
	}

	paths := []string{
		filepath.Join(m.checkDir, m.encodeHeight(seed.Height())),
		filepath.Join(m.valDir, m.encodeHash(seed.Header.ValidatorsHash)),
	}
	for _, p := range paths {
		outf, err := os.Create(p)
		if os.IsExist(err) {
			continue
		}
		if err == nil {
			var n int
			wire.WriteBinary(seed, outf, &n, &err)
			// we can't use defer as we are in a loop... make sure we close
			outf.Close()
		}
		// unknown error in creating or writing immediately breaks
		if err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

func (m Provider) GetByHeight(h int) (certifiers.Seed, error) {
	// first we look for exact match, then search...
	s := certifiers.Seed{}
	path := filepath.Join(m.checkDir, m.encodeHeight(h))
	inf, err := os.Open(path)
	// if no match, make a best attempt to find one lower
	if os.IsNotExist(err) {
		inf, err = m.searchForHeight(h)
	}

	if err == nil {
		defer inf.Close()
		var n int
		wire.ReadBinaryPtr(&s, inf, 0, &n, &err)
	}

	// error here on read file or parse file
	return s, errors.WithStack(err)
}

// search for height, looks for a file with highest height < h
// return certifiers.ErrSeedNotFound() if not there...
func (m Provider) searchForHeight(h int) (*os.File, error) {
	d, err := os.Open(m.checkDir)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	files, err := d.Readdirnames(0)

	d.Close()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	desired := m.encodeHeight(h)
	sort.Strings(files)
	i := sort.SearchStrings(files, desired)
	if i == 0 {
		return nil, certifiers.ErrSeedNotFound()
	}
	found := files[i-1]

	path := filepath.Join(m.checkDir, found)
	inf, err := os.Open(path)
	return inf, errors.WithStack(err)
}

func (m Provider) GetByHash(hash []byte) (certifiers.Seed, error) {
	s := certifiers.Seed{}
	path := filepath.Join(m.valDir, m.encodeHash(hash))
	inf, err := os.Open(path)
	if os.IsNotExist(err) {
		return s, certifiers.ErrSeedNotFound()
	}

	if err == nil {
		defer inf.Close()
		var n int
		wire.ReadBinaryPtr(&s, inf, 0, &n, &err)
	}

	// error here on read file or parse file
	return s, errors.WithStack(err)
}
