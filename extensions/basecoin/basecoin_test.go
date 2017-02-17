package basecoin_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	bc "github.com/tendermint/basecoin/types"
	wire "github.com/tendermint/go-wire"
	"github.com/tendermint/light-client/cryptostore"
	"github.com/tendermint/light-client/extensions/basecoin"
	"github.com/tendermint/light-client/rpc/tests"
	"github.com/tendermint/light-client/storage/memstorage"
)

func TestBasecoinSendTx(t *testing.T) {
	assert, require := assert.New(t), require.New(t)
	// node must parse basecoin values
	n := tests.GetNode()
	n.ValueReader = basecoin.BasecoinValues{}

	// store the keys somewhere
	keys := cryptostore.New(
		cryptostore.GenEd25519,
		cryptostore.SecretBox,
		memstorage.New(),
	)

	// make a user
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

	// set some data with SetOption
	acct := bc.Account{
		PubKey: k.PubKey,
		Balance: bc.Coins{
			{Denom: "ATOM", Amount: 1000},
			{Denom: "ETH", Amount: 150},
		},
	}
	acctjson := string(wire.JSONBytes(&acct))
	log := bcapp.SetOption("base/account", acctjson)
	assert.Equal("Success", log)

	// wait for one more block, so this data is commited and in block
	n.WaitForHeight(q.Height + 1)

	// try querying node for this info - some data
	q2, err := n.Query("/account", addr)
	require.Nil(err)
	require.NotNil(q2.Value)
	// we should read an account back
	qa, ok := q2.Value.(basecoin.Account)
	require.True(ok, "%#v", q2.Value)
	// and make sure it looks write
	assert.Equal(basecoin.AccountType, qa.Type)
	qav := qa.Value
	assert.Equal(acct.Balance, qav.Balance)
	assert.Equal(acct.PubKey, qav.PubKey.PubKey)
}
