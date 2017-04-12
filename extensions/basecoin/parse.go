package basecoin

import (
	"encoding/json"

	"github.com/pkg/errors"
	bc "github.com/tendermint/basecoin/types"
)

/***
 I removed most of the custom json parsing.
 We still need a better way to multi-plex over app data in basecoin
***/

func parseAppTx(data []byte, appData AppDataReader) (*bc.AppTx, error) {
	var tx txApp
	err := json.Unmarshal(data, &tx)
	if err != nil {
		return nil, errors.Wrap(err, "parse apptx")
	}
	atx, err := tx.toBasecoin(appData)
	return &atx, err
}

type txApp struct {
	Gas     int64           `json:"gas"`   // Gas
	Fee     bc.Coin         `json:"fee"`   // Fee
	Name    string          `json:"name"`  // Which plugin
	Input   bc.TxInput      `json:"input"` // Hmmm do we want coins?
	Type    string          `json:"type"`  // which tx type for this plugin
	AppData json.RawMessage `json:"appdata"`
}

func (t txApp) toBasecoin(appData AppDataReader) (bc.AppTx, error) {
	data, err := appData(t.Name, t.Type, t.AppData)
	return bc.AppTx{
		Gas:   t.Gas,
		Fee:   t.Fee,
		Name:  t.Name,
		Input: t.Input,
		Data:  data,
	}, err
}
