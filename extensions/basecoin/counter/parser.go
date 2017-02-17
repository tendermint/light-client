package counter

import (
	"encoding/json"

	"github.com/tendermint/basecoin/plugins/counter"
	wire "github.com/tendermint/go-wire"
)

func ReadCounterTx(data []byte) ([]byte, error) {
	var tx counter.CounterTx
	err := json.Unmarshal(data, &tx)
	if err != nil {
		return nil, err
	}
	return wire.BinaryBytes(tx), nil
}
