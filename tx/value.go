package tx

import (
	"encoding/hex"
	"encoding/json"

	b58 "github.com/jbenet/go-base58"
	lc "github.com/tendermint/light-client"
)

// HexData let's us treat a byte slice as hex data, rather than default base64
type HexData []byte

func (h *HexData) UnmarshalJSON(b []byte) (err error) {
	var s string
	err = json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	// and interpret that string as hex
	bin, err := hex.DecodeString(s)
	*h = bin
	return err
}

func (h HexData) MarshalJSON() ([]byte, error) {
	s := hex.EncodeToString(h)
	return json.Marshal(s)
}

// B58Data let's us treat a byte slice as base58, like bitcoin addresses
type B58Data []byte

func (d *B58Data) UnmarshalJSON(b []byte) (err error) {
	var s string
	err = json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	// TODO: modify code to return error (tendermint fork?)
	bin := b58.Decode(s)
	*d = bin
	return err
}

func (d B58Data) MarshalJSON() ([]byte, error) {
	s := b58.Encode(d)
	return json.Marshal(s)
}

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
