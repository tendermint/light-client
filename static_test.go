package lightclient_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	crypto "github.com/tendermint/go-crypto"
	lc "github.com/tendermint/light-client"
	"github.com/tendermint/tendermint/rpc/client"
	"github.com/tendermint/tendermint/types"
)

func getLocalClient() client.Local {
	return client.NewLocal(node)
}

// TODO: run this without a real tendermint node...
// Test in certifiers????
func TestNodeAuditing(t *testing.T) {
	assert, require := assert.New(t), require.New(t)

	cl := getLocalClient()
	time.Sleep(100 * time.Millisecond)

	// get initial validators
	vals, err := cl.Validators()
	require.Nil(err, "%+v", err)

	// let's grab a header and make sure it is proper
	commit, err := cl.Commit(vals.BlockHeight)
	require.Nil(err, "%+v", err)
	check := lc.NewCheckpoint(commit)

	// let's see if this checkpoint will validate our new validator set now
	vals, err = cl.Validators()
	require.Nil(err, "%+v", err)
	err = check.CheckValidators(vals.Validators)
	assert.Nil(err, "%+v", err)

	// make an invalid set, which should fail
	pk := crypto.GenPrivKeyEd25519().PubKey()
	badval := append(vals.Validators, types.NewValidator(pk, 3))
	err = check.CheckValidators(badval)
	assert.NotNil(err)
}
