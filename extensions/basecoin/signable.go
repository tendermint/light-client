package basecoin

import (
	"encoding/json"

	"github.com/pkg/errors"
	lc "github.com/tendermint/light-client"
)

type TxType struct {
	Type string
	Data json.RawMessage
}

type BasecoinTx struct {
	ChainID string
}

func (t BasecoinTx) assertSignableReader() lc.SignableReader {
	return t
}

// Turn json into a signable object
func (t BasecoinTx) ReadSignable(data []byte) (lc.Signable, error) {
	var tx TxType
	err := json.Unmarshal(data, &tx)
	if err != nil {
		return nil, errors.Wrap(err, "Read JSON Tx")
	}
	// switch to tx type
	if tx.Type == "sendtx" {
		return t.readSendTx(tx.Data)
	} else if tx.Type == "apptx" {
		return t.readAppTx(tx.Data)
	}
	return nil, errors.Errorf("Unknown type: %s", tx.Type)
}
