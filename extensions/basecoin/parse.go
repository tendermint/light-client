package basecoin

import (
	"encoding/json"

	"github.com/pkg/errors"
	bc "github.com/tendermint/basecoin/types"
	"github.com/tendermint/light-client/tx"
)

/**** TODO: all this ugliness must go away when we refactor json parsing ***/

func parseSendTx(data []byte) (*bc.SendTx, error) {
	var tx txSend
	err := json.Unmarshal(data, &tx)
	if err != nil {
		return nil, errors.Wrap(err, "parse sendtx")
	}
	stx := tx.toBasecoin()
	return &stx, nil
}

func parseAppTx(data []byte) (*bc.AppTx, error) {
	var tx txApp
	err := json.Unmarshal(data, &tx)
	if err != nil {
		return nil, errors.Wrap(err, "parse apptx")
	}
	atx := tx.toBasecoin()
	return &atx, nil
}

// WARNING/NOTE: does not handle serializing sigs, as we don't take them over json
type txInput struct {
	Address  tx.HexData    `json:"address"`  // Hash of the PubKey
	Coins    bc.Coins      `json:"coins"`    //
	Sequence int           `json:"sequence"` // Must be 1 greater than the last committed TxInput
	PubKey   tx.JSONPubKey `json:"pub_key"`  // Is present iff Sequence == 1
}

func (t txInput) toBasecoin() bc.TxInput {
	return bc.TxInput{
		Address:  t.Address,
		Coins:    t.Coins,
		Sequence: t.Sequence,
		PubKey:   t.PubKey.PubKey,
	}
}

type txOutput struct {
	Address tx.HexData `json:"address"` // Hash of the PubKey
	Coins   bc.Coins   `json:"coins"`   //
}

func (t txOutput) toBasecoin() bc.TxOutput {
	return bc.TxOutput{
		Address: t.Address,
		Coins:   t.Coins,
	}
}

type txSend struct {
	Gas     int64      `json:"gas"` // Gas
	Fee     bc.Coin    `json:"fee"` // Fee
	Inputs  []txInput  `json:"inputs"`
	Outputs []txOutput `json:"outputs"`
}

func (t txSend) toBasecoin() bc.SendTx {
	ins := make([]bc.TxInput, len(t.Inputs))
	for i := range t.Inputs {
		ins[i] = t.Inputs[i].toBasecoin()
	}

	outs := make([]bc.TxOutput, len(t.Outputs))
	for i := range t.Outputs {
		outs[i] = t.Outputs[i].toBasecoin()
	}

	return bc.SendTx{
		Gas:     t.Gas,
		Fee:     t.Fee,
		Inputs:  ins,
		Outputs: outs,
	}
}

type txApp struct {
	Gas   int64      `json:"gas"`   // Gas
	Fee   bc.Coin    `json:"fee"`   // Fee
	Name  string     `json:"type"`  // Which plugin
	Input txInput    `json:"input"` // Hmmm do we want coins?
	Data  tx.HexData `json:"data"`
}

func (t txApp) toBasecoin() bc.AppTx {
	return bc.AppTx{
		Gas:   t.Gas,
		Fee:   t.Fee,
		Name:  t.Name,
		Input: t.Input.toBasecoin(),
		Data:  t.Data,
	}
}
