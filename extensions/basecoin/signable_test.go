package basecoin

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	bc "github.com/tendermint/basecoin/types"
	crypto "github.com/tendermint/go-crypto"
)

func TestSendTxJSON(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	// notice pubkey is one string not that funky array
	raw := []byte(`{
    "type": "sendtx",
    "data": {
      "gas": 22,
      "fee": {"denom": "ETH", "amount": 1},
      "inputs": [{
        "address": "4d8908785ec867139ca02e71a717c01fa506b96a",
        "coins": [{"denom": "ETH", "amount": 21}],
        "sequence": 1,
        "pub_key": "01d7fb176319af0c126c4c4c7851cf7c66340e7df8763f0aa9700ebae32a955401"
      }],
      "outputs": [{
        "address": "9f31a3ac6b1468402aac5593ae9e82a041457f5f",
        "coins": [{"denom": "ETH", "amount": 20}]
      }]
    }
  }`)
	sr := NewBasecoinTx("foo")
	sig, err := sr.ReadSignable(raw)
	require.Nil(err)
	stx, ok := sig.(*SendTx)
	require.True(ok)

	tx := stx.tx
	require.NotNil(tx)
	assert.EqualValues(22, tx.Gas)
	assert.Equal("ETH", tx.Fee.Denom)
	if assert.Equal(1, len(tx.Inputs)) {
		validateInput(t, tx.Inputs[0])
	}
	if assert.Equal(1, len(tx.Outputs)) {
		out := tx.Outputs[0]
		addr, err := hex.DecodeString("9f31a3ac6b1468402aac5593ae9e82a041457f5f")
		require.Nil(err)
		assert.EqualValues(addr, out.Address)
		assert.Equal(1, len(out.Coins))
		assert.EqualValues(20, out.Coins[0].Amount)
		assert.EqualValues("ETH", out.Coins[0].Denom)
	}
}

func validateInput(t *testing.T, in bc.TxInput) {
	assert := assert.New(t)
	require := require.New(t)
	addr, err := hex.DecodeString("4d8908785ec867139ca02e71a717c01fa506b96a")
	require.Nil(err)
	assert.Equal(addr, in.Address)
	assert.Equal(1, len(in.Coins))
	assert.EqualValues(21, in.Coins[0].Amount)
	require.NotNil(in.PubKey)
	// ensure type byte reflected proper
	pk, ok := in.PubKey.(crypto.PubKeyEd25519)
	assert.True(ok)
	// check the first byte - d7 - decoded proper
	assert.Equal(pk[0], byte(215))
}

type demoData struct {
	Key   string
	Value string
}

func demoParse(data []byte) ([]byte, error) {
	var d demoData
	err := json.Unmarshal(data, &d)
	res := fmt.Sprintf("%s=%s", d.Key, d.Value)
	return []byte(res), err
}

func TestAppTxJSON(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	// notice pubkey is one string not that funky array
	raw := []byte(`{
    "type": "apptx",
    "data": {
      "name": "demo",
      "gas": 78,
      "fee": {"denom": "ATOM", "amount": 5},
      "input": {
        "address": "4d8908785ec867139ca02e71a717c01fa506b96a",
        "coins": [{"denom": "ATOM", "amount": 21}],
        "sequence": 1,
        "pub_key": "01d7fb176319af0c126c4c4c7851cf7c66340e7df8763f0aa9700ebae32a955401"
      },
      "type": "create",
      "appdata": {
        "key": "hello",
        "value": "bonjour"
      }
    }
  }`)
	sr := NewBasecoinTx("foo")
	// note: we must register all tx types we wish to support
	sr.RegisterParser("demo", "create", demoParse)

	sig, err := sr.ReadSignable(raw)
	require.Nil(err)
	atx, ok := sig.(*AppTx)
	require.True(ok)

	tx := atx.tx
	require.NotNil(tx)
	assert.EqualValues(78, tx.Gas)
	assert.Equal("ATOM", tx.Fee.Denom)
	assert.EqualValues(5, tx.Fee.Amount)
	assert.Equal("demo", tx.Name)

	// verify the input as above
	validateInput(t, tx.Input)

	// and make sure out special app data is properly formated
	assert.Equal("hello=bonjour", string(tx.Data))
}
