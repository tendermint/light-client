package basecoin

import lc "github.com/tendermint/light-client"

type BasecoinValues struct{}

// Turn merkle binary into a json-able struct
func (t BasecoinValues) ReadValue(key, value []byte) (lc.Value, error) {
	// TODO
	return nil, nil
}

func (v BasecoinValues) assertValueReader() lc.ValueReader {
	return v
}
