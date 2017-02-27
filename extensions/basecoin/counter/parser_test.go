package counter

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tendermint/basecoin/plugins/counter"
	"github.com/tendermint/basecoin/types"
	wire "github.com/tendermint/go-wire"
)

func TestParse(t *testing.T) {
	// convert data to binary
	require := require.New(t)
	data := `{"valid": true, "fee": [{"denom": "wtf", "amount": 123}]}`
	bin, err := ReadCounterTx([]byte(data))
	require.Nil(err, "%+v", err)

	// recover the original
	var tx counter.CounterTx
	err = wire.ReadBinaryBytes(bin, &tx)
	require.Nil(err, "%+v", err)

	// verify it is correct
	require.True(tx.Valid)
	require.Equal(1, len(tx.Fee))
	require.Equal("wtf", tx.Fee[0].Denom)
	require.EqualValues(123, tx.Fee[0].Amount)
}

func TestValue(t *testing.T) {
	require := require.New(t)
	state := counter.CounterPluginState{
		Counter:   15,
		TotalFees: types.Coins{{Denom: "change", Amount: 420}},
	}
	val := wire.BinaryBytes(state)

	out, err := Value{}.ReadValue(nil, val)
	require.Nil(err)
	require.Equal(val, out.Bytes())
	cnt, ok := out.(Counter)
	require.True(ok)
	require.Equal(state, cnt.CounterPluginState)
}
