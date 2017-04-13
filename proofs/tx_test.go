package proofs_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/light-client/proofs"
	merktest "github.com/tendermint/merkleeyes/testutil"
)

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
	h := br.Height

	// let's get a proof for our tx
	pr, err := prover.Get(tx, uint64(h))
	require.Nil(err, "%+v", err)

	// make sure bad queries return errors
	_, err = prover.Get([]byte("no-such-tx"), uint64(h))
	require.NotNil(err)
	_, err = prover.Get(tx, uint64(h+1))
	require.NotNil(err)

	// matches and validates with post-tx header
	check := getCheckForHeight(t, cl, h)
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
