package basecoin_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/basecoin/app"
	bc "github.com/tendermint/basecoin/types"
	data "github.com/tendermint/go-data"
	keys "github.com/tendermint/go-keys"
	"github.com/tendermint/go-keys/cryptostore"
	"github.com/tendermint/go-keys/storage/memstorage"
	"github.com/tendermint/light-client/extensions/basecoin"
	"github.com/tendermint/light-client/extensions/basecoin/counter"
	"github.com/tendermint/light-client/rpc/tests"
)

const DefaultAlgo = "ed25519"

func makeUser(t *testing.T, keys cryptostore.Manager, name, pass string) keys.Info {
	k, err := keys.Create(name, pass, DefaultAlgo)
	require.Nil(t, err)
	return k
}

func setAcct(t *testing.T, bcapp *app.Basecoin, acct *bc.Account) {
	acctjson, err := data.ToJSON(acct)
	require.Nil(t, err, "%+v", err)
	log := bcapp.SetOption("base/account", string(acctjson))
	require.Equal(t, "Success", log)
}

// TestBasecoinSetOption tests whether we can create an account using
// SetOption (as per genesis), and read it back properly from the db
func TestBasecoinSetOption(t *testing.T) {
	assert, require := assert.New(t), require.New(t)
	// node must parse basecoin values
	n := tests.GetNode()
	n.ValueReader = basecoin.NewBasecoinValues()

	// store the keys somewhere
	keys := cryptostore.New(
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
	n.ValueReader = basecoin.NewBasecoinValues()

	// store the keys somewhere
	keys := cryptostore.New(
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
	sr.RegisterParser("counter", "counter", counter.ReadCounterTx)
	key_data, err := data.ToJSON(k1.PubKey)
	require.Nil(err)
	raw := fmt.Sprintf(`{
    "type": "sendtx",
    "data": {
      "gas": 22,
      "fee": {"denom": "ATOM", "amount": 1},
      "inputs": [{
        "address": "%X",
        "coins": [{"denom": "ATOM", "amount": 21}],
        "sequence": 1,
        "pub_key": %s
      }],
      "outputs": [{
        "address": "%X",
        "coins": [{"denom": "ATOM", "amount": 20}]
      }]
    }
  }`, addr1, key_data, addr2)
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

// TestBasecoinAppTx executes an AppTx on the counter app
func TestBasecoinAppTx(t *testing.T) {
	assert, require := assert.New(t), require.New(t)
	// node must parse basecoin values
	n := tests.GetNode()
	val := basecoin.NewBasecoinValues()
	val.RegisterPlugin(counter.Value{})
	// TODO: register plugin parser
	n.ValueReader = val

	// register the plugin for the sender
	sr := basecoin.NewBasecoinTx(ChainID)
	sr.RegisterParser("counter", "counter", counter.ReadCounterTx)

	// store the keys somewhere
	keys := cryptostore.New(
		cryptostore.SecretBox,
		memstorage.New(),
	)
	name, pass := "connie", "a-count-ant"
	k := makeUser(t, keys, name, pass)
	addr := k.PubKey.Address()

	// set some data with SetOption
	acct := bc.Account{
		PubKey:  k.PubKey,
		Balance: bc.Coins{{Denom: "gold", Amount: 5432}},
	}
	setAcct(t, bcapp, &acct)

	key_data, err := data.ToJSON(k.PubKey)
	require.Nil(err)
	// now, let's generate a tx
	raw := fmt.Sprintf(`{
    "type": "apptx",
    "data": {
      "name": "counter",
      "gas": 22,
      "fee": {
        "denom": "gold",
        "amount": 2
      },
      "input": {
        "address": "%X",
        "coins": [{
          "denom": "gold",
          "amount": 20
        }],
        "sequence": 1,
        "pub_key": %s
      },
      "type": "counter",
      "appdata": {
        "valid": true,
        "fee": [{
          "denom": "gold",
          "amount": 5
        }]
      }
    }
  }`, addr, key_data)
	sig, err := sr.ReadSignable([]byte(raw))
	require.Nil(err, "%+v", err)
	_, ok := sig.(*basecoin.AppTx)
	require.True(ok)

	// sign and send it
	keys.Sign(name, pass, sig)
	tx, err := sig.TxBytes()
	require.Nil(err)
	bres, err := n.Broadcast(tx)
	require.Nil(err, "%+v", err)
	require.True(bres.IsOk(), "%#v", bres)

	// wait for one more block...
	q, err := n.Query("/account", addr)
	require.Nil(err)
	n.WaitForHeight(q.Height + 1)

	// and the both fees were deducted
	q, err = n.Query("/account", addr)
	require.Nil(err)
	require.NotNil(q.Value)
	qav := q.Value.(basecoin.Account).Value
	// TODO: fix counter, currently we lose all input, even if not
	// used up by fees
	assert.Equal(bc.Coins{{Denom: "gold", Amount: 5412}}, qav.Balance)

	// query counter state!
	cntkey := []byte("CounterPlugin.State")
	cq, err := n.Query("/key", cntkey)
	require.Nil(err)
	require.NotNil(cq.Value)
	// make sure it's parsed
	cstate, ok := cq.Value.(counter.Counter)
	require.True(ok)
	require.Equal(1, cstate.Counter)
	require.Equal(bc.Coins{{Denom: "gold", Amount: 5}}, cstate.TotalFees)

	// and make sure it is nice json
	_, err = json.Marshal(cq.Value)
	require.Nil(err)
}
