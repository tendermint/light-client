package certifiers_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/light-client/certifiers"
)

func TestInquirerValidPath(t *testing.T) {
	assert, require := assert.New(t), require.New(t)
	trust := certifiers.NewMemStoreProvider()
	source := certifiers.NewMemStoreProvider()

	// set up the validators to generate test blocks
	var vote int64 = 10
	keys := certifiers.GenValKeys(5)
	vals := keys.ToValidators(vote, 0)

	// initialize a certifier with the initial state
	chainID := "inquiry-test"
	cert := certifiers.NewInquiring(chainID, vals, trust, source)

	// construct a bunch of seeds, each with one more height than the last
	count := 50
	seeds := make([]certifiers.Seed, count)
	for i := 0; i < count; i++ {
		// extend the keys by 1 each time
		keys = keys.Extend(1)
		vals = keys.ToValidators(vote, 0)
		h := 20 + 10*i
		appHash := []byte(fmt.Sprintf("h=%d", h))
		cp := keys.GenCheckpoint(chainID, h, nil, vals, appHash, 0, len(keys))
		seeds[i] = certifiers.Seed{cp, vals}
	}

	check := seeds[count-1].Checkpoint

	// this should fail validation....
	err := cert.Certify(check)
	require.NotNil(err)

	// add a few seed in the middle should be insufficient
	for i := 10; i < 13; i++ {
		err := cert.SeedSource.StoreSeed(seeds[i])
		require.Nil(err)
	}
	err = cert.Certify(check)
	assert.NotNil(err)

	// with more info, we succeed
	for i := 0; i < count; i++ {
		err := cert.SeedSource.StoreSeed(seeds[i])
		require.Nil(err)
	}
	err = cert.Certify(check)
	assert.Nil(err, "%+v", err)
}

func TestInquirerMinimalPath(t *testing.T) {
	assert, require := assert.New(t), require.New(t)
	trust := certifiers.NewMemStoreProvider()
	source := certifiers.NewMemStoreProvider()

	// set up the validators to generate test blocks
	var vote int64 = 10
	keys := certifiers.GenValKeys(5)
	vals := keys.ToValidators(vote, 0)

	// initialize a certifier with the initial state
	chainID := "minimal-path"
	cert := certifiers.NewInquiring(chainID, vals, trust, source)

	// construct a bunch of seeds, each with one more height than the last
	count := 12
	seeds := make([]certifiers.Seed, count)
	for i := 0; i < count; i++ {
		// extend the validators, so we are just below 2/3
		keys = keys.Extend(len(keys)/2 - 1)
		vals = keys.ToValidators(vote, 0)
		h := 5 + 10*i
		appHash := []byte(fmt.Sprintf("h=%d", h))
		cp := keys.GenCheckpoint(chainID, h, nil, vals, appHash, 0, len(keys))
		seeds[i] = certifiers.Seed{cp, vals}
	}
	check := seeds[count-1].Checkpoint

	// this should fail validation....
	err := cert.Certify(check)
	require.NotNil(err)

	// add a few seed in the middle should be insufficient
	for i := 5; i < 8; i++ {
		err := cert.SeedSource.StoreSeed(seeds[i])
		require.Nil(err)
	}
	err = cert.Certify(check)
	assert.NotNil(err)

	// with more info, we succeed
	for i := 0; i < count; i++ {
		err := cert.SeedSource.StoreSeed(seeds[i])
		require.Nil(err)
	}
	err = cert.Certify(check)
	assert.Nil(err, "%+v", err)
}
