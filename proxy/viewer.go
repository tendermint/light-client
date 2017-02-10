package proxy

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	lc "github.com/tendermint/light-client"
	"github.com/tendermint/light-client/proxy/types"
	"github.com/tendermint/light-client/util"
)

type Viewer struct {
	lc.Checker
	lc.Searcher
	lc.ValueReader
	util.Auditor
}

func NewViewer(checker lc.Checker, searcher lc.Searcher,
	cert lc.Certifier) Viewer {
	return Viewer{
		Checker:  checker,
		Searcher: searcher,
		Auditor:  util.NewAuditor(cert),
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

	res, err := v.Query(path, data)
	if err != nil {
		writeError(w, err)
		return
	}

	if !res.Code.IsOK() {
		writeCode(w, renderQueryFail(res), 400)
	}

	resp, err := renderQuery(res, false)
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
	p, err := v.Prove(key)
	if err != nil {
		writeError(w, err)
		return
	}
	if !p.Code.IsOK() {
		writeCode(w, renderQueryFail(p), 400)
	}

	// waiting until the signatures are ready (if needed)
	err = v.WaitForHeight(p.Height)
	if err != nil {
		writeError(w, err)
		return
	}

	// get the next block height
	// TODO: when there is no height?
	block, err := v.SignedHeader(p.Height)
	if err != nil {
		writeError(w, err)
		return
	}
	err = v.Audit(p.Key, p.Value.Bytes(), p.Proof, block)
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

func renderQuery(r lc.TmQueryResult, proven bool) (*types.QueryResponse, error) {
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

func renderQueryFail(r lc.TmQueryResult) *types.GenericResponse {
	return &types.GenericResponse{
		Code: r.Code,
		Log:  r.Log,
	}
}
