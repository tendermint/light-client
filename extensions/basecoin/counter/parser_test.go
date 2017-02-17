package counter

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tendermint/basecoin/plugins/counter"
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
