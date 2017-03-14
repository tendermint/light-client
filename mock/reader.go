package mock

import (
	keys "github.com/tendermint/go-keys"
	wire "github.com/tendermint/go-wire"
	lc "github.com/tendermint/light-client"
)

const (
	typeOneSig   = byte(0x81)
	typeMultiSig = byte(0x82)
)

type wrapper struct{ keys.Signable }

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

func (r reader) ReadSignable(data []byte) (keys.Signable, error) {
	res := wrapper{}
	err := wire.ReadBinaryBytes(data, &res)
	return res.Signable, err
}
