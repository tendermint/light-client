package files

import (
	"encoding/hex"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"

	"github.com/pkg/errors"

	"github.com/tendermint/light-client/certifiers"
	certerr "github.com/tendermint/light-client/certifiers/errors"
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
// headers/*Commits, and thus store it separately. This leaves us
// with three issues:
//
// 1. Given a validator hash, retrieve the validator set if previously stored
// 2. Given a block height, find the *Commit with the highest height <= h
// 3. Given a Seed, store it quickly to satisfy 1 and 2
//
// Note that we do not worry about caching, as that can be achieved by
// pairing this with a MemStoreProvider and CacheProvider from certifiers
type Provider struct {
	valDir   string
	checkDir string
}

var _ certifiers.Provider = &Provider{}

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
	err := seed.ValidateBasic(seed.Commit.Header.ChainID)
	if err != nil {
		return err
	}

	paths := []string{
		filepath.Join(m.checkDir, m.encodeHeight(seed.Height())),
		filepath.Join(m.valDir, m.encodeHash(seed.Header.ValidatorsHash)),
	}
	for _, p := range paths {
		err := seed.Write(p)
		// unknown error in creating or writing immediately breaks
		if err != nil {
			return err
		}
	}
	return nil
}

func (m Provider) GetByHeight(h int) (certifiers.Seed, error) {
	// first we look for exact match, then search...
	path := filepath.Join(m.checkDir, m.encodeHeight(h))
	seed, err := certifiers.LoadSeed(path)
	if certerr.IsSeedNotFoundErr(err) {
		path, err = m.searchForHeight(h)
		if err == nil {
			seed, err = certifiers.LoadSeed(path)
		}
	}
	return seed, err
}

func (p *Provider) LatestSeed() (seed certifiers.Seed, err error) {
	return p.GetByHeight(math.MaxInt32 - 1)
}

// search for height, looks for a file with highest height < h
// return certifiers.ErrSeedNotFound() if not there...
func (m Provider) searchForHeight(h int) (string, error) {
	d, err := os.Open(m.checkDir)
	if err != nil {
		return "", errors.WithStack(err)
	}
	files, err := d.Readdirnames(0)

	d.Close()
	if err != nil {
		return "", errors.WithStack(err)
	}

	desired := m.encodeHeight(h)
	sort.Strings(files)
	i := sort.SearchStrings(files, desired)
	if i == 0 {
		return "", certerr.ErrSeedNotFound()
	}
	found := files[i-1]
	path := filepath.Join(m.checkDir, found)
	return path, errors.WithStack(err)
}

func (m Provider) GetByHash(hash []byte) (certifiers.Seed, error) {
	path := filepath.Join(m.valDir, m.encodeHash(hash))
	return certifiers.LoadSeed(path)
}
