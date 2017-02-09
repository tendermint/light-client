package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	crypto "github.com/tendermint/go-crypto"
	lc "github.com/tendermint/light-client"
	"github.com/tendermint/light-client/rpc"
	"github.com/tendermint/light-client/util"
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
	assert.Equal(v, qr.Value.Bytes())
	assert.Nil(qr.Proof)

	// and we get some proof, we can even decipher
	pr, err := n.Prove(k)
	require.Nil(err)
	assert.True(pr.Code.IsOK())
	assert.Equal(k, pr.Key)
	assert.Equal(v, pr.Value.Bytes())
	assert.NotNil(pr.Proof)

	p := pr.Proof
	root := p.Root()
	assert.NotNil(root)
	// this proof validates our data
	assert.True(p.Verify(k, v, root))
	// but not some mixed-up data
	assert.False(p.Verify(v, k, root))
}

func TestNodeHeaders(t *testing.T) {
	assert, require := assert.New(t), require.New(t)
	n := GetNode()

	// send some data
	_, _, tx := TestTxKV()
	br, err := n.Broadcast(tx)
	require.Nil(err)
	require.True(br.Code.IsOK())

	// get the validator set
	vals, err := n.Validators()
	require.Nil(err)
	assert.Equal(1, len(vals.Validators))

	// get a signed header
	height := vals.BlockHeight - 1 // FIXME: looking for sig queries height+1...
	block, err := n.SignedHeader(height)
	require.Nil(err, "%+v", err)
	assert.Equal(height, block.Header.Height)
	assert.Equal(vals.BlockHeight, block.LastHeight)
	assert.Equal(1, len(block.Votes))

	// try to certify this header is proper
	cert := rpc.StaticCertifier{Vals: vals.Validators}
	err = cert.Certify(block)
	assert.Nil(err, "%+v", err)

	// but with other validators, not so good...
	newKey := crypto.GenPrivKeySecp256k1()
	pk := newKey.PubKey()
	// not enough to block... yet
	power := (vals.Validators[0].VotingPower / 2) - 1
	nv := append(vals.Validators, lc.TmValidator{
		Address:     pk.Address(),
		PubKey:      pk,
		VotingPower: power,
	})
	badcert := rpc.StaticCertifier{Vals: nv}
	err = badcert.Certify(block)
	assert.Nil(err, "%+v", err)

	// but let's give this fake validator just a bit more power....
	// and we no longer have quorum
	nv[1].VotingPower += 2
	err = badcert.Certify(block)
	assert.NotNil(err)
}

// TestNodeAuditing attempts to go the whole route to validate data in
// the system.  It needs to get a proof, get the headers, validate all sigs
// and be happy.  Also need to test failures....
func TestNodeAuditing(t *testing.T) {
	assert, require := assert.New(t), require.New(t)
	n := GetNode()

	// send some data
	k, v, tx := TestTxKV()
	br, err := n.Broadcast(tx)
	require.Nil(err)
	require.True(br.Code.IsOK())

	// let's query for this data now
	pr, err := n.Prove(k)
	require.Nil(err)
	assert.True(pr.Code.IsOK())
	assert.Equal(k, pr.Key)
	assert.Equal(v, pr.Value.Bytes())
	assert.NotNil(pr.Proof)
	proot := pr.Proof.Root() // the roothash from the proof

	// get the validator set (also find current height)
	vals, err := n.Validators()
	require.Nil(err)
	cert := rpc.StaticCertifier{Vals: vals.Validators}
	auditor := util.NewAuditor(cert)
	height := vals.BlockHeight

	// we need to push some more blocks on here, so we can query...
	// this whole need to wait one-two blocks to get a proof
	k2, v2, tx2 := TestTxKV()
	_, err = n.Broadcast(tx2)
	require.Nil(err)
	pr2, err := n.Prove(k2)
	require.Nil(err)
	assert.NotNil(pr2.Proof)

	// get a signed header of height+1 which should have apphash for height
	oldblock, err := n.SignedHeader(height)
	require.Nil(err, "%+v", err)
	block, err := n.SignedHeader(height + 1)
	require.Nil(err, "%+v", err)
	// let's see if the root hash matches the proof
	assert.Equal(proot, block.Header.AppHash)

	// okay, now let's do a full audit...
	err = auditor.Audit(k, v, pr.Proof, block)
	require.Nil(err, "%+v", err)
	// will fail for the wrong block header... or wrong values... or wrong proof
	err = auditor.Audit(k, v, pr.Proof, oldblock)
	require.NotNil(err)
	err = auditor.Audit(k, v, pr2.Proof, oldblock)
	require.NotNil(err)
	err = auditor.Audit(k2, v2, pr.Proof, block)
	require.NotNil(err)
	err = auditor.Audit(k2, v2, pr2.Proof, block)
	require.NotNil(err)

	// wait, let's try to prove the new one as well
	_, err = n.SignedHeader(height + 2)
	require.NotNil(err)

	// oops... we have to move the block along first... rrr....
	for i := 0; i < 2; i++ {
		_, _, dtx := TestTxKV()
		dbr, err := n.Broadcast(dtx)
		require.Nil(err)
		require.True(dbr.IsOk())
	}

	// now we can prove the new entry as well
	// wtf?  +3?? I thought it would be +2?  If only the dummy app returned
	// the height that relates to the proof...
	block2, err := n.SignedHeader(height + 3)
	require.Nil(err, "%+v", err)
	err = auditor.Audit(k2, v2, pr2.Proof, block2)
	require.Nil(err, "%+v", err)

}
