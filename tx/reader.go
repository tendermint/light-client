package tx

import (
	keys "github.com/tendermint/go-keys"
	wire "github.com/tendermint/go-wire"
)

const (
	typeOneSig   = byte(0x01)
	typeMultiSig = byte(0x02)
)

type wrapper struct{ keys.Signable }

func init() {
	wire.RegisterInterface(
		wrapper{},
		wire.ConcreteType{&OneSig{}, typeOneSig},
		wire.ConcreteType{&MultiSig{}, typeMultiSig},
	)
}

func ReadSignableBinary(data []byte) (keys.Signable, error) {
	res := wrapper{}
	err := wire.ReadBinaryBytes(data, &res)
	return res.Signable, err
}
