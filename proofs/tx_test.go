package proofs_test

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/light-client/proofs"
	merktest "github.com/tendermint/merkleeyes/testutil"
	"github.com/tendermint/tendermint/rpc/client"
)

// findTx returns the firts block height with at least one tx
func findTx(cl client.Client, start, end int) (uint64, error) {
	// fmt.Printf("Searching %d to %d\n", end, start)
	headers, err := cl.BlockchainInfo(start, end)
	if err != nil {
		return 0, err
	}
	for _, head := range headers.BlockMetas {
		h := head.Header.Height
		dh := head.Header.DataHash
		// fmt.Printf("%d: %X\n", h, dh)
		if len(dh) > 0 {
			return uint64(h), nil
		}
	}
	return 0, errors.New("No tx found")
}

func TestTxProofs(t *testing.T) {
	assert, require := assert.New(t), require.New(t)

	cl := getLocalClient()
	prover := proofs.NewTxProver(cl)
	time.Sleep(200 * time.Millisecond)

	precheck := getCurrentCheck(t, cl)

	// great, let's store some data here, and make more checks....
	_, _, tx := merktest.MakeTxKV()
	br, err := cl.BroadcastTxCommit(tx)
	require.Nil(err, "%+v", err)
	require.EqualValues(0, br.CheckTx.GetCode())
	require.EqualValues(0, br.DeliverTx.GetCode())

	h, err := findTx(cl, precheck.Height()-1, precheck.Height()+20)
	require.Nil(err, "%+v", err)

	// unfortunately we cannot tell the server to give us any height
	// other than the most recent, so 0 is the only choice :(
	pr, err := prover.Get(tx, h)
	require.Nil(err, "%+v", err)
	check := getCheckForHeight(t, cl, int(h))

	// matches and validates with post-tx header
	err = pr.Validate(check)
	assert.Nil(err, "%+v", err)

	// doesn't matches with pre-tx header
	err = pr.Validate(precheck)
	assert.NotNil(err)

	// make sure it has the values we want
	txpr, ok := pr.(proofs.TxProof)
	if assert.True(ok) {
		assert.EqualValues(tx, txpr.Tx())
	}

	// make sure we read/write properly, and any changes to the serialized
	// object are invalid proof (2000 random attempts)
	testSerialization(t, prover, pr, check, 2000)
}

// // validate all tx in the block
// block, err := cl.Block(check.Height())
// require.Nil(err, "%+v", err)
// err = check.CheckTxs(block.Block.Data.Txs)
// assert.Nil(err, "%+v", err)

// oh, i would like the know the hieght of the broadcast_commit.....
// so i could verify that tx :(
