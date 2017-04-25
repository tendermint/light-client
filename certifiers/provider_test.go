package certifiers_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/light-client/certifiers"
)

func TestMemProvider(t *testing.T) {
	p := certifiers.NewMemStoreProvider()
	checkProvider(t, p, "test-mem", "empty")
}

func TestCacheProvider(t *testing.T) {
	p := certifiers.NewCacheProvider(
		certifiers.NewMissingProvider(),
		certifiers.NewMemStoreProvider(),
		certifiers.NewMissingProvider(),
	)
	checkProvider(t, p, "test-cache", "kjfhekfhkewhgit")
}

func checkProvider(t *testing.T, p certifiers.Provider, chainID, app string) {
	assert, require := assert.New(t), require.New(t)
	appHash := []byte(app)
	keys := certifiers.GenValKeys(5)
	count := 10

	// make a bunch of seeds...
	seeds := make([]certifiers.Seed, count)
	for i := 0; i < count; i++ {
		// two seeds for each validator, to check how we handle dups
		// (10, 0), (10, 1), (10, 1), (10, 2), (10, 2), ...
		vals := keys.ToValidators(10, int64(count/2))
		h := 20 + 10*i
		check := keys.GenCheckpoint(chainID, h, nil, vals, appHash, 0, 5)
		seeds[i] = certifiers.Seed{check, vals}
	}

	// check provider is empty
	seed, err := p.GetByHeight(20)
	require.NotNil(err)
	assert.True(certifiers.IsSeedNotFoundErr(err))

	seed, err = p.GetByHash(seeds[3].Hash())
	require.NotNil(err)
	assert.True(certifiers.IsSeedNotFoundErr(err))

	// now add them all to the provider
	for _, s := range seeds {
		err = p.StoreSeed(s)
		require.Nil(err)
		// and make sure we can get it back
		s2, err := p.GetByHash(s.Hash())
		assert.Nil(err)
		assert.Equal(s, s2)
		// by height as well
		s2, err = p.GetByHeight(s.Height())
		assert.Nil(err)
		assert.Equal(s, s2)
	}

	// make sure we get the last hash if we overstep
	seed, err = p.GetByHeight(5000)
	if assert.Nil(err) {
		assert.Equal(seeds[count-1].Height(), seed.Height())
		assert.Equal(seeds[count-1], seed)
	}

	// and middle ones as well
	seed, err = p.GetByHeight(47)
	if assert.Nil(err) {
		// we only step by 10, so 40 must be the one below this
		assert.Equal(40, seed.Height())
	}

}

// this will make a get height, and if it is good, set the data as well
func checkGetHeight(t *testing.T, p certifiers.Provider, ask, expect int) {
	seed, err := p.GetByHeight(ask)
	require.Nil(t, err, "%+v", err)
	if assert.Equal(t, expect, seed.Height()) {
		err = p.StoreSeed(seed)
		require.Nil(t, err, "%+v", err)
	}
}

func TestCacheGetsBestHeight(t *testing.T) {
	// assert, require := assert.New(t), require.New(t)
	require := require.New(t)

	// we will write data to the second level of the cache (p2),
	// and see what gets cached, stored in
	p := certifiers.NewMemStoreProvider()
	p2 := certifiers.NewMemStoreProvider()
	cp := certifiers.NewCacheProvider(p, p2)

	chainID := "cache-best-height"
	appHash := []byte("01234567")
	keys := certifiers.GenValKeys(5)
	count := 10

	// set a bunch of seeds
	for i := 0; i < count; i++ {
		vals := keys.ToValidators(10, int64(count/2))
		h := 10 * (i + 1)
		check := keys.GenCheckpoint(chainID, h, nil, vals, appHash, 0, 5)
		seed := certifiers.Seed{check, vals}
		err := p2.StoreSeed(seed)
		require.Nil(err)
	}

	// let's get a few heights from the cache and set them proper
	checkGetHeight(t, cp, 57, 50)
	checkGetHeight(t, cp, 33, 30)

	// make sure they are set in p as well (but nothing else)
	checkGetHeight(t, p, 44, 30)
	checkGetHeight(t, p, 50, 50)
	checkGetHeight(t, p, 99, 50)

	// now, query the cache for a higher value
	checkGetHeight(t, p2, 99, 90)
	checkGetHeight(t, cp, 99, 90)
}
