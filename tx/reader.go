package tx

import (
	wire "github.com/tendermint/go-wire"
	lc "github.com/tendermint/light-client"
)

const (
	typeOneSig   = byte(0x01)
	typeMultiSig = byte(0x02)
)

type wrapper struct{ lc.Signable }

func init() {
	wire.RegisterInterface(
		wrapper{},
		wire.ConcreteType{&OneSig{}, typeOneSig},
		wire.ConcreteType{&MultiSig{}, typeMultiSig},
	)
}

type reader struct{}

// Reader constructs a SignableReader that can parse OneSig and MultiSig
//
// TODO: add some args to configure go-wire, rather than relying on init???7
func Reader() lc.SignableReader {
	return reader{}
}

func (r reader) ReadSignable(data []byte) (lc.Signable, error) {
	res := wrapper{}
	err := wire.ReadBinaryBytes(data, &res)
	return res.Signable, err
}
