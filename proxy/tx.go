package proxy

import (
	"net/http"

	"github.com/gorilla/mux"
	keys "github.com/tendermint/go-keys"
	lc "github.com/tendermint/light-client"
	"github.com/tendermint/light-client/proxy/types"
	"github.com/tendermint/tendermint/rpc/client"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
)

type TxSigner struct {
	lc.SignableReader
	lc.Poster
}

func NewTxSigner(server client.ABCIClient, signer keys.Signer,
	reader lc.SignableReader) TxSigner {

	return TxSigner{
		SignableReader: reader,
		Poster:         lc.NewPoster(server, signer),
	}
}

func (t TxSigner) PostTransaction(w http.ResponseWriter, r *http.Request) {
	req := types.PostTxRequest{}
	err := readRequest(r, &req)
	if err != nil {
		writeError(w, err)
		return
	}

	tx, err := t.ReadSignable(req.Data)
	if err != nil {
		writeError(w, err)
		return
	}

	res, err := t.Post(tx, req.Name, req.Passphrase)
	if err != nil {
		writeError(w, err)
		return
	}

	resp := renderBroadcast(res)
	writeSuccess(w, &resp)
}

func (t TxSigner) Register(r *mux.Router) {
	r.HandleFunc("/", t.PostTransaction).Methods("POST")
}

func renderBroadcast(r *ctypes.ResultBroadcastTxCommit) types.GenericResponse {
	return types.GenericResponse{
		Code: int32(r.DeliverTx.Code),
		Data: r.DeliverTx.Data,
		Log:  r.DeliverTx.Log,
	}
}
