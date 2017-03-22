package proxy

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	abci "github.com/tendermint/abci/types"
	lc "github.com/tendermint/light-client"
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

	resp, err := renderQuery(q, false)
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

	// get the proof
	pr, err := v.ABCIQuery("/key", key, true)
	if err != nil {
		writeError(w, err)
		return
	}
	p := pr.Response
	if p.Code != 0 {
		writeCode(w, renderQueryFail(p), 400)
	}

	// waiting until the signatures are ready (if needed)
	h := int(p.Height)
	err = client.WaitForHeight(v, h, nil)
	if err != nil {
		writeError(w, err)
		return
	}

	// get the next block height
	// TODO: when there is no height?
	com, err := v.Commit(h)
	if err != nil {
		writeError(w, err)
		return
	}
	check := lc.NewCheckpoint(com)

	// let's see if this is valid
	err = v.Certify(check)
	if err != nil {
		writeError(w, err)
		return
	}

	// now use the checkpoint to validate proof
	proof, err := lc.MerkleReader.ReadProof(p.Proof)
	if err != nil {
		writeError(w, err)
		return
	}
	err = check.CheckAppState(p.Key, p.Value, proof)
	if err != nil {
		writeError(w, err)
		return
	}

	// D00d, we are good!
	// that was one hard won boolean flag
	resp, err := renderQuery(p, true)
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

func renderQuery(r abci.ResponseQuery, proven bool) (*types.QueryResponse, error) {
	value, err := json.Marshal(r.Value)
	if err != nil {
		return nil, err
	}
	return &types.QueryResponse{
		Height: r.Height,
		Key:    r.Key,
		Value:  value,
		Proven: proven,
	}, nil
}

func renderQueryFail(r abci.ResponseQuery) *types.GenericResponse {
	return &types.GenericResponse{
		Code: int32(r.Code),
		Log:  r.Log,
	}
}
