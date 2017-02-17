package basecoin_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/light-client/cryptostore"
	"github.com/tendermint/light-client/rpc/tests"
	"github.com/tendermint/light-client/storage/memstorage"
)

func TestBasecoinSendTx(t *testing.T) {
	assert, require := assert.New(t), require.New(t)
	n := tests.GetNode()
	keys := cryptostore.New(
		cryptostore.GenEd25519,
		cryptostore.SecretBox,
		memstorage.New(),
	)

	name, pass := "freddy", "**mercury**"
	err := keys.Create(name, pass)
	require.Nil(err)
	k, err := keys.Get(name)
	require.Nil(err)

	addr := k.PubKey.Address()

	// try querying node for this info - empty
	q, err := n.Query("/account", addr)
	require.Nil(err)
	assert.Equal(q.Value, nil)

	// let's set some options here...
	// bcapp.SetOption("base/account")

	// wait for one more block

	// try querying node for this info - empty

}
