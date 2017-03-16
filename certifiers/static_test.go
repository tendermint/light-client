package certifiers_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	lc "github.com/tendermint/light-client"
	"github.com/tendermint/light-client/certifiers"
	"github.com/tendermint/tendermint/rpc/client"
	rpctest "github.com/tendermint/tendermint/rpc/test"
)

func getLocalClient() client.Local {
	return client.NewLocal(node)
}

func TestStaticCert(t *testing.T) {
	require := require.New(t)
	// assert, require := assert.New(t), require.New(t)
	cl := getLocalClient()
	chainID := rpctest.GetConfig().GetString("chain_id")

	// set up the certifier by getting the validator set
	// (in real code we don't trust a node this naively)
	vals, err := cl.Validators()
	require.Nil(err, "%+v", err)
	cert := certifiers.NewStatic(chainID, vals.Validators)

	h := vals.BlockHeight + 1
	client.WaitForHeight(cl, h, nil)
	commit, err := cl.Commit(h)
	require.Nil(err, "%+v", err)

	check := lc.NewCheckpoint(commit)
	err = cert.Certify(check)
	require.Nil(err, "%+v", err)
}
