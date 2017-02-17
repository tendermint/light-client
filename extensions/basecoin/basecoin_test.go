package basecoin_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/basecoin/app"
	bc "github.com/tendermint/basecoin/types"
	wire "github.com/tendermint/go-wire"
	lc "github.com/tendermint/light-client"
	"github.com/tendermint/light-client/cryptostore"
	"github.com/tendermint/light-client/extensions/basecoin"
	"github.com/tendermint/light-client/rpc/tests"
	"github.com/tendermint/light-client/storage/memstorage"
)

func makeUser(t *testing.T, keys cryptostore.Manager, name, pass string) lc.KeyInfo {
	err := keys.Create(name, pass)
	require.Nil(t, err)
	k, err := keys.Get(name)
	require.Nil(t, err)
	return k
}

func setAcct(t *testing.T, bcapp *app.Basecoin, acct *bc.Account) {
	acctjson := string(wire.JSONBytes(acct))
	log := bcapp.SetOption("base/account", acctjson)
	require.Equal(t, "Success", log)
}

// TestBasecoinSetOption tests whether we can create an account using
// SetOption (as per genesis), and read it back properly from the db
func TestBasecoinSetOption(t *testing.T) {
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
	k := makeUser(t, keys, name, pass)
	addr := k.PubKey.Address()

	// try querying node for this info - empty
	q, err := n.Query("/account", addr)
	require.Nil(err)
	assert.Nil(q.Value)

	// set some data with SetOption
	acct := bc.Account{
		PubKey: k.PubKey,
		Balance: bc.Coins{
			{Denom: "ATOM", Amount: 1000},
			{Denom: "ETH", Amount: 150},
		},
	}
	setAcct(t, bcapp, &acct)

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

// TestBasecoinSendTx sets up an account and send money to a second
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

	// make two users user
	n1, p1 := "donny", "donny.darko"
	n2, p2 := "fuzzy", "was.a.bear!"
	k1 := makeUser(t, keys, n1, p1)
	k2 := makeUser(t, keys, n2, p2)
	addr1 := k1.PubKey.Address()
	addr2 := k2.PubKey.Address()

	// set some data with SetOption
	acct := bc.Account{
		PubKey:  k1.PubKey,
		Balance: bc.Coins{{Denom: "ATOM", Amount: 234}},
	}
	setAcct(t, bcapp, &acct)

	// check there is no data in addr2
	q, err := n.Query("/account", addr2)
	require.Nil(err)
	assert.Nil(q.Value)

	// now, let's generate a tx
	sr := basecoin.NewBasecoinTx(ChainID)
	raw := fmt.Sprintf(`{
    "type": "sendtx",
    "data": {
      "gas": 22,
      "fee": {"denom": "ATOM", "amount": 1},
      "inputs": [{
        "address": "%X",
        "coins": [{"denom": "ATOM", "amount": 21}],
        "sequence": 1,
        "pub_key": "%X"
      }],
      "outputs": [{
        "address": "%X",
        "coins": [{"denom": "ATOM", "amount": 20}]
      }]
    }
  }`, addr1, k1.PubKey.Bytes(), addr2)
	sig, err := sr.ReadSignable([]byte(raw))
	require.Nil(err)
	_, ok := sig.(*basecoin.SendTx)
	require.True(ok)

	// send it
	tx, err := sig.TxBytes()
	require.Nil(err)
	bres, err := n.Broadcast(tx)
	require.Nil(err, "%+v", err)
	require.False(bres.IsOk(), "%#v", bres)

	// but sign it first
	keys.Sign(n1, p1, sig)
	tx, err = sig.TxBytes()
	require.Nil(err)
	bres, err = n.Broadcast(tx)
	require.Nil(err, "%+v", err)
	require.True(bres.IsOk(), "%#v", bres)

	// wait for one more block...
	q, err = n.Query("/account", addr2)
	require.Nil(err)
	n.WaitForHeight(q.Height + 1)

	// make sure the money arrived
	q, err = n.Query("/account", addr2)
	require.Nil(err)
	require.NotNil(q.Value)
	qav := q.Value.(basecoin.Account).Value
	assert.Equal(bc.Coins{{Denom: "ATOM", Amount: 20}}, qav.Balance)

	// and was deducted
	q, err = n.Query("/account", addr1)
	require.Nil(err)
	require.NotNil(q.Value)
	qav = q.Value.(basecoin.Account).Value
	assert.Equal(bc.Coins{{Denom: "ATOM", Amount: 213}}, qav.Balance)
}
