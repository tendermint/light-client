package counter

import (
	"bytes"
	"encoding/json"
	"errors"

	"github.com/tendermint/basecoin/plugins/counter"
	wire "github.com/tendermint/go-wire"
	lc "github.com/tendermint/light-client"
)

func ReadCounterTx(data []byte) ([]byte, error) {
	var tx counter.CounterTx
	err := json.Unmarshal(data, &tx)
	if err != nil {
		return nil, err
	}
	return wire.BinaryBytes(tx), nil
}

type Value struct{}

func (_ Value) ReadValue(key, value []byte) (lc.Value, error) {
	target := []byte("CounterPlugin.State")
	if len(key) == 0 || bytes.Equal(target, key) {
		var cpState counter.CounterPluginState
		err := wire.ReadBinaryBytes(value, &cpState)
		if err == nil {
			return Counter{cpState}, nil
		}
	}
	return nil, errors.New("Cannot parse counter")
}

type Counter struct {
	counter.CounterPluginState
}

func (c Counter) Bytes() []byte {
	return wire.BinaryBytes(c.CounterPluginState)
}
