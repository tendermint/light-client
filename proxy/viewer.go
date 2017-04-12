package proxy

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	abci "github.com/tendermint/abci/types"
	lc "github.com/tendermint/light-client"
	"github.com/tendermint/light-client/proofs"
	"github.com/tendermint/light-client/proxy/types"
	"github.com/tendermint/tendermint/rpc/client"
)

type Viewer struct {
	client.Client
	// lc.ValueReader
	lc.Certifier
}

func NewViewer(cl client.Client, cert lc.Certifier) Viewer {
	return Viewer{
		Client:    cl,
		Certifier: cert,
	}
}

func (v Viewer) QueryData(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	path, hexdata := vars["path"], vars["data"]

	// decode using hex
	data, err := hex.DecodeString(hexdata)
	if err != nil {
		writeError(w, errors.New("data must be hexidecimal"))
		return
	}

	res, err := v.ABCIQuery("/"+path, data, false)
	if err != nil {
		writeError(w, err)
		return
	}

	q := res.Response
	if q.Code != 0 {
		writeCode(w, renderQueryFail(q), 400)
		return
	}

	resp, err := renderQuery(q)
	if err == nil {
		writeSuccess(w, resp)
	} else {
		writeError(w, err)
	}
}

func (v Viewer) ProveData(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hexkey := vars["key"]

	// decode using hex
	key, err := hex.DecodeString(hexkey)
	if err != nil {
		writeError(w, errors.New("key must be hexidecimal"))
		return
	}

	prover := proofs.NewAppProver(v)
	pr, err := prover.Get(key, 0)
	if err != nil {
		writeError(w, err)
		return
	}

	// waiting until the signatures are ready (if needed)
	h := int(pr.BlockHeight())
	err = client.WaitForHeight(v, h, nil)
	if err != nil {
		writeError(w, err)
		return
	}

	// get the next block height
	com, err := v.Commit(h)
	if err != nil {
		writeError(w, err)
		return
	}
	check := lc.NewCheckpoint(com)

	// let's see if the checkpoint is properly signed
	err = v.Certify(check)
	if err != nil {
		writeError(w, err)
		return
	}

	// now make sure the proof matches our new header
	err = pr.Validate(check)
	if err != nil {
		writeError(w, err)
		return
	}

	// D00d, we are good!
	// that was one hard won boolean flag
	resp, err := renderProof(pr.(proofs.AppProof))
	if err == nil {
		writeSuccess(w, resp)
	} else {
		writeError(w, err)
	}
}

func (v Viewer) Register(r *mux.Router) {
	r.HandleFunc("/query/{path}/{data}", v.QueryData).Methods("GET")
	r.HandleFunc("/proof/{key}", v.ProveData).Methods("GET")
}

func renderQuery(r abci.ResponseQuery) (*types.QueryResponse, error) {
	value, err := json.Marshal(r.Value)
	if err != nil {
		return nil, err
	}
	return &types.QueryResponse{
		Height: r.Height,
		Key:    r.Key,
		Value:  value,
		Proven: false,
	}, nil
}

func renderProof(p proofs.AppProof) (*types.QueryResponse, error) {
	value, err := json.Marshal(p.Value)
	if err != nil {
		return nil, err
	}
	return &types.QueryResponse{
		Height: p.Height,
		Key:    p.Key,
		Value:  value,
		Proven: true,
	}, nil
}

func renderQueryFail(r abci.ResponseQuery) *types.GenericResponse {
	return &types.GenericResponse{
		Code: int32(r.Code),
		Log:  r.Log,
	}
}
