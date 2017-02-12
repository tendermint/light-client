package basecoin

import (
	lc "github.com/tendermint/light-client"
	"github.com/tendermint/light-client/tx"
)

type BasecoinValues struct{}

// Turn merkle binary into a json-able struct
func (t BasecoinValues) ReadValue(key, value []byte) (lc.Value, error) {
	// TODO - something more than hex
	return tx.NewValue(value), nil
}

func (v BasecoinValues) assertValueReader() lc.ValueReader {
	return v
}
