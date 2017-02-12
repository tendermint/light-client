package tx

import (
	"encoding/hex"
	"errors"

	lc "github.com/tendermint/light-client"
)

type HexData []byte

func (h *HexData) UnmarshalJSON(b []byte) (err error) {
	l := len(b)
	if l < 2 || b[0] != '"' || b[l-1] != '"' {
		return errors.New("hex string must be enclosed with quotes")
	}
	bin, err := hex.DecodeString(string(b[1 : l-1]))
	*h = bin
	return err
}

func (h HexData) MarshalJSON() ([]byte, error) {
	hex := `"` + hex.EncodeToString(h) + `"`
	return []byte(hex), nil
}

// TODO: Marshal/Unmarshal json as hex

const RawValueType = "raw"

type RawValue struct {
	Type  string  `json:"type"`
	Value HexData `json:"value"`
}

func NewValue(val []byte) RawValue {
	return RawValue{
		Type:  RawValueType,
		Value: HexData(val),
	}
}

func (v RawValue) Bytes() []byte {
	return []byte(v.Value)
}

func (v RawValue) assertValue() lc.Value {
	return v
}
