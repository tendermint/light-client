package tx

import (
	crypto "github.com/tendermint/go-crypto"
	data "github.com/tendermint/go-data"
	keys "github.com/tendermint/go-keys"
)

const (
	typeOneSig   = byte(0x01)
	typeMultiSig = byte(0x02)
	nameOneSig   = "sig"
	nameMultiSig = "multi"
)

var _ keys.Signable = Sig{}
var TxMapper data.Mapper

func init() {
	TxMapper = data.NewMapper(Sig{}).
		RegisterInterface(&OneSig{}, nameOneSig, typeOneSig).
		RegisterInterface(&MultiSig{}, nameMultiSig, typeMultiSig)
}

type SigInner interface {
	SignBytes() []byte
	Sign(pubkey crypto.PubKey, sig crypto.Signature) error
	Signers() ([]crypto.PubKey, error)
}

// Sig is what is exported, and handles serialization
type Sig struct {
	SigInner
}

func (s Sig) TxBytes() ([]byte, error) {
	return data.ToWire(s)
}

// func (s Sig) ReadSignable(data []byte) (keys.Signable, error) {
// 	data, err := TxMapper.FromJSON(data)
// 	return data.(keys.Signable), err
// }
