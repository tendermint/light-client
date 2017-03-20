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
	)
	checkProvider(t, p, "test-cache", "kjfhekfhkewh")
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
	assert.True(certifiers.SeedNotFound(err))

	seed, err = p.GetByHash(seeds[3].Hash())
	require.NotNil(err)
	assert.True(certifiers.SeedNotFound(err))

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
