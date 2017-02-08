package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/go-merkle"
)

func TestNodeQuery(t *testing.T) {
	assert, require := assert.New(t), require.New(t)
	n := GetNode()

	// send it works
	k, v, tx := TestTxKV()
	br, err := n.Broadcast(tx)
	require.Nil(err)
	require.True(br.Code.IsOK())

	// query it is there
	qr, err := n.Query("/key", k)
	require.Nil(err)
	assert.True(qr.Code.IsOK())
	assert.Equal(v, qr.Value)
	assert.Nil(qr.Proof)

	// and we get some proof, we can even decipher
	pr, err := n.Prove(k)
	require.Nil(err)
	assert.True(pr.Code.IsOK())
	assert.Equal(k, pr.Key)
	assert.Equal(v, pr.Value)
	assert.NotNil(pr.Proof)

	proof, err := merkle.ReadProof(pr.Proof)
	require.Nil(err)
	root := proof.RootHash
	assert.NotNil(root)
	// this proof validates our data
	assert.True(proof.Verify(k, v, root))
	// but not some mixed-up data
	assert.False(proof.Verify(v, k, root))
}

func TestNodeHeaders(t *testing.T) {
	assert, require := assert.New(t), require.New(t)
	n := GetNode()

	// get the validator set
	vals, err := n.Validators()
	require.Nil(err)
	assert.Equal(1, len(vals))

	// send some data
	_, _, tx := TestTxKV()
	br, err := n.Broadcast(tx)
	require.Nil(err)
	require.True(br.Code.IsOK())

	// get a signed header
	height := uint64(1) // TODO - better
	head, err := n.SignedHeader(height)
	require.Nil(err, "%+v", err)
	assert.Equal(height, head.Header.Height)
	assert.Equal(1, len(head.Votes))
}
