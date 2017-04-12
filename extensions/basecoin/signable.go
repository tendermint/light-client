package basecoin

import (
	"encoding/json"

	"github.com/pkg/errors"
	keys "github.com/tendermint/go-keys"
)

type TxType struct {
	Type string           `json:"type"`
	Data *json.RawMessage `json:"data"`
}

type BasecoinTx struct {
	chainID string
	appMap  map[string]TxReader
}

func NewBasecoinTx(chainID string) BasecoinTx {
	return BasecoinTx{
		chainID: chainID,
		appMap:  map[string]TxReader{},
	}
}

// AppDataReader takes a plugin name and txType and some json
// and serializes it into a binary format
type AppDataReader func(name, txType string, json []byte) ([]byte, error)

// TxReader handles parsing and serializing one particular type
type TxReader func(json []byte) ([]byte, error)

// Turn json into a signable object
func (t BasecoinTx) ReadSignable(data []byte) (keys.Signable, error) {
	var tx TxType
	err := json.Unmarshal(data, &tx)
	if err != nil {
		return nil, errors.Wrap(err, "Read JSON Tx")
	}
	// switch to tx type
	if tx.Type == "sendtx" {
		return t.readSendTx(*tx.Data)
	} else if tx.Type == "apptx" {
		return t.readAppTx(*tx.Data)
	}
	return nil, errors.Errorf("Unknown type: %s", tx.Type)
}

func (t BasecoinTx) RegisterParser(name, txType string, reader TxReader) {
	key := name + "/" + txType
	t.appMap[key] = reader
}

func (t BasecoinTx) appData(name, txType string, json []byte) ([]byte, error) {
	key := name + "/" + txType
	reader, ok := t.appMap[key]
	if !ok {
		return nil, errors.Errorf("No registered parser for %s/%s", name, txType)
	}
	return reader(json)
}
