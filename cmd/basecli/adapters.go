package main

import (
	"encoding/hex"

	btypes "github.com/tendermint/basecoin/types"
	wire "github.com/tendermint/go-wire"
	"github.com/tendermint/light-client/proofs"
)

type AccountPresenter struct{}

func (_ AccountPresenter) MakeKey(str string) ([]byte, error) {
	res, err := hex.DecodeString(str)
	if err == nil {
		res = append([]byte("base/a/"), res...)
	}
	return res, err
}

func (_ AccountPresenter) ParseData(raw []byte) (interface{}, error) {
	var acc *btypes.Account
	err := wire.ReadBinaryBytes(raw, &acc)
	return acc, err
}

type BaseTxPresenter struct {
	proofs.RawPresenter // this handles MakeKey as hex bytes
}

func (_ BaseTxPresenter) ParseData(raw []byte) (interface{}, error) {
	var tx btypes.TxS
	err := wire.ReadBinaryBytes(raw, &tx)
	return tx, err
}
