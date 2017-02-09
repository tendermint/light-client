package mock

import lc "github.com/tendermint/light-client"

// ByteValue is a simple way to pass byte slices unparsed as Values
// meant only for test cases, as usually you will want the ValueReader
// to actually make sense of the structure
type ByteValue []byte

func (b ByteValue) Bytes() []byte {
	return []byte(b)
}

func (b ByteValue) assertValue() lc.Value {
	return b
}

// ValueReader returns a mock ValueReader for test cases
func ValueReader() lc.ValueReader {
	return ByteValueReader{}
}

// ByteValueReader is a simple implementation that just wraps the bytes
// in ByteValue.
//
// Intended for testing where there is no app-specific data structure
type ByteValueReader struct{}

func (b ByteValueReader) ReadValue(key, value []byte) (lc.Value, error) {
	return ByteValue(value), nil
}
