package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	btypes "github.com/tendermint/basecoin/types"
	wire "github.com/tendermint/go-wire"
	lightclient "github.com/tendermint/light-client"
	"github.com/tendermint/light-client/commands"
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

// SendTXReader allows us to create SendTx
type SendTxReader struct {
	ChainID string
}

func (t SendTxReader) ReadTxJSON(data []byte) (interface{}, error) {
	var tx btypes.SendTx
	err := json.Unmarshal(data, &tx)
	send := SendTx{
		chainID: t.ChainID,
		Tx:      &tx,
	}
	return &send, errors.Wrap(err, "parse sendtx")
}

type SendTxMaker struct{}

func (m SendTxMaker) MakeReader() (lightclient.TxReader, error) {
	chainID := viper.GetString(commands.ChainFlag)
	return SendTxReader{ChainID: chainID}, nil
}

type SendFlags struct {
	To       string
	From     string
	Amount   string
	Fee      string
	Gas      int
	Sequence int
}

func (m SendTxMaker) Flags() (*flag.FlagSet, interface{}) {
	fs := flag.NewFlagSet("foobar", flag.ContinueOnError)
	fs.String("to", "", "Destination address for the bits")
	fs.String("from", "", "Sender address for the tx")
	fs.String("amount", "", "Coins to send in the format <amt><coin>,<amt><coin>...")
	fs.String("fee", "", "Coins for the transaction fee of the format <amt><coin>")
	fs.Int("gas", 0, "Amount of gas for this transaction")
	fs.Int("sequence", -1, "Sequence number for this transaction")
	return fs, &SendFlags{}
}

func (t SendTxReader) ReadTxFlags(flags interface{}) (interface{}, error) {
	data := flags.(*SendFlags)
	fmt.Printf("Data: %#v\n", data)
	return nil, errors.New("gotcha")
}
