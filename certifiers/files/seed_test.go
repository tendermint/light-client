package files

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	cmn "github.com/tendermint/tmlibs/common"

	"github.com/tendermint/light-client/certifiers"
)

func tmpFile() string {
	suffix := cmn.RandStr(16)
	return filepath.Join(os.TempDir(), "seed-test-"+suffix)
}

func TestSerializeSeeds(t *testing.T) {
	assert, require := assert.New(t), require.New(t)

	// some constants
	appHash := []byte("some crazy thing")
	chainID := "ser-ial"
	h := 25

	// build a seed
	keys := certifiers.GenValKeys(5)
	vals := keys.ToValidators(10, 0)
	check := keys.GenCommit(chainID, h, nil, vals, appHash, 0, 5)
	seed := certifiers.Seed{check, vals}

	require.Equal(h, seed.Height())
	require.Equal(vals.Hash(), seed.ValidatorsHash())

	// try read/write with json
	jfile := tmpFile()
	defer os.Remove(jfile)
	jseed, err := LoadSeedJSON(jfile)
	assert.NotNil(err)
	err = WriteSeedJSON(seed, jfile)
	require.Nil(err)
	jseed, err = LoadSeedJSON(jfile)
	assert.Nil(err, "%+v", err)
	assert.Equal(h, jseed.Height())
	assert.Equal(vals.Hash(), jseed.ValidatorsHash())

	// try read/write with binary
	bfile := tmpFile()
	defer os.Remove(bfile)
	bseed, err := LoadSeed(bfile)
	assert.NotNil(err)
	err = WriteSeed(seed, bfile)
	require.Nil(err)
	bseed, err = LoadSeed(bfile)
	assert.Nil(err, "%+v", err)
	assert.Equal(h, bseed.Height())
	assert.Equal(vals.Hash(), bseed.ValidatorsHash())

	// make sure they don't read the other format (different)
	_, err = LoadSeed(jfile)
	assert.NotNil(err)
	_, err = LoadSeedJSON(bfile)
	assert.NotNil(err)
}
