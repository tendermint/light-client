package client_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/light-client/certifiers"
	"github.com/tendermint/light-client/certifiers/client"
	rpctest "github.com/tendermint/tendermint/rpc/test"
)

func TestProvider(t *testing.T) {
	assert, require := assert.New(t), require.New(t)

	cfg := rpctest.GetConfig()
	rpcAddr := cfg.RPC.ListenAddress
	chainID := cfg.ChainID
	p := client.NewHTTP(rpcAddr)
	require.NotNil(t, p)

	// let it produce some blocks
	time.Sleep(500 * time.Millisecond)

	// let's get the highest block
	seed, err := p.GetByHeight(5000)
	require.Nil(err, "%+v", err)
	sh := seed.Height()
	vhash := seed.Header.ValidatorsHash
	assert.True(sh < 5000)

	// let's check this is valid somehow
	assert.Nil(seed.ValidateBasic(chainID))
	cert := certifiers.NewStatic(chainID, seed.Validators)

	// can't get a lower one
	seed, err = p.GetByHeight(sh - 1)
	assert.NotNil(err)
	assert.True(certifiers.IsSeedNotFoundErr(err))

	// also get by hash (given the match)
	seed, err = p.GetByHash(vhash)
	require.Nil(err, "%+v", err)
	require.Equal(vhash, seed.Header.ValidatorsHash)
	err = cert.Certify(seed.Checkpoint)
	assert.Nil(err, "%+v", err)

	// get by hash fails without match
	seed, err = p.GetByHash([]byte("foobar"))
	assert.NotNil(err)
	assert.True(certifiers.IsSeedNotFoundErr(err))

	// storing the seed silently ignored
	err = p.StoreSeed(seed)
	assert.Nil(err, "%+v", err)
}
