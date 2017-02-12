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

func ReadSignableBinary(data []byte) (lc.Signable, error) {
	res := wrapper{}
	err := wire.ReadBinaryBytes(data, &res)
	return res.Signable, err
}
