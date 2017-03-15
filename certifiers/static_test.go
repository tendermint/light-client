package certifiers_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	crypto "github.com/tendermint/go-crypto"
	lc "github.com/tendermint/light-client"
	"github.com/tendermint/light-client/certifiers"
	merktest "github.com/tendermint/merkleeyes/testutil"
	"github.com/tendermint/tendermint/rpc/client"
	rpctest "github.com/tendermint/tendermint/rpc/test"
	"github.com/tendermint/tendermint/types"
)

func getLocalClient() client.Local {
	return client.NewLocal(node)
}

// TestNodeAuditing attempts to go the whole route to validate data in
// the system.  It needs to get a proof, get the headers, validate all sigs
// and be happy.  Also need to test failures....
func TestNodeAuditing(t *testing.T) {
	assert, require := assert.New(t), require.New(t)

	cl := getLocalClient()
	chainID := rpctest.GetConfig().GetString("chain_id")

	// set up the certifier by getting the validator set
	// (in real code we don't trust a node this naively)
	vals, err := cl.Validators()
	require.Nil(err, "%+v", err)
	cert := certifiers.NewStatic(chainID, vals.Validators)

	// great, let's store some data here, and make more checks....
	k, v, tx := merktest.MakeTxKV()
	br, err := cl.BroadcastTxCommit(tx)
	require.Nil(err, "%+v", err)
	require.EqualValues(0, br.CheckTx.GetCode())
	require.EqualValues(0, br.DeliverTx.GetCode())

	// let's grab a header and make sure it is proper
	// oh, i wish the broadcast would give us a height...
	commit, err := cl.Commit(vals.BlockHeight + 1)
	require.Nil(err, "%+v", err)
	check := lc.NewCheckpoint(commit)
	err = cert.Certify(check)
	require.Nil(err, "%+v", err)

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

	// let's query for this data now
	pres, err := cl.ABCIQuery("/store", k, true)
	require.Nil(err, "%+v", err)
	pr := pres.Response
	assert.Equal(k, pr.Key)
	assert.Equal(v, pr.Value)
	assert.NotNil(pr.Proof)
	proof, err := certifiers.MerkleReader.ReadProof(pr.Proof)
	require.Nil(err, "%+v", err)

	// and get a new checkpoint to match
	h := int(pr.Height)
	client.WaitForHeight(cl, h, nil)
	commit2, err := cl.Commit(h)
	require.Nil(err, "%+v", err)
	check2 := lc.NewCheckpoint(commit2)
	err = cert.Certify(check2)
	require.Nil(err, "%+v", err)

	// validate the proof
	err = check2.CheckAppState(pr.Key, pr.Value, proof)
	assert.Nil(err, "%+v", err)

	// other header doesn't validate
	if check.Height() != check2.Height() {
		err = check.CheckAppState(pr.Key, pr.Value, proof)
		assert.NotNil(err, "%+v", err)
	}

	// and an erroneous proof doesn't validate
	pr.Key[5] = 0
	err = check2.CheckAppState(pr.Key, pr.Value, proof)
	assert.NotNil(err)

	// validate all tx in the block
	block, err := cl.Block(check.Height())
	require.Nil(err, "%+v", err)
	err = check.CheckTxs(block.Block.Data.Txs)
	assert.Nil(err, "%+v", err)

	// oh, i would like the know the hieght of the broadcast_commit.....
	// so i could verify that tx :(
}
