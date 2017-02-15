package basecoin

import (
	"encoding/json"

	"github.com/pkg/errors"
	bc "github.com/tendermint/basecoin/types"
)

/**** TODO: all this ugliness must go away when we refactor json parsing ***/

func parseSendTx(data []byte) (*bc.SendTx, error) {
	var tx bc.SendTx
	err := json.Unmarshal(data, &tx)
	return &tx, errors.Wrap(err, "parse sendtx")
}

func parseAppTx(data []byte) (*bc.AppTx, error) {
	var tx bc.AppTx
	err := json.Unmarshal(data, &tx)
	return &tx, errors.Wrap(err, "parse apptx")
}
