package proxy

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	lc "github.com/tendermint/light-client"
	"github.com/tendermint/light-client/proxy/types"
)

type KeyServer struct {
	manager lc.KeyManager
}

func NewKeyServer(manager lc.KeyManager) KeyServer {
	return KeyServer{
		manager: manager,
	}
}

func (k KeyServer) GenerateKey(w http.ResponseWriter, r *http.Request) {
	req := types.CreateKeyRequest{}
	err := readRequest(r, &req)
	if err != nil {
		writeError(w, err)
		return
	}

	err = k.manager.Create(req.Name, req.Passphrase)
	if err != nil {
		writeError(w, err)
		return
	}

	key, err := k.manager.Get(req.Name)
	if err != nil {
		writeError(w, err)
		return
	}

	resp := renderKey(key)
	writeSuccess(w, &resp)
}

func (k KeyServer) GetKey(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	key, err := k.manager.Get(name)
	if err != nil {
		writeError(w, err)
		return
	}

	resp := renderKey(key)
	writeSuccess(w, &resp)
}

func (k KeyServer) ListKeys(w http.ResponseWriter, r *http.Request) {

	keys, err := k.manager.List()
	if err != nil {
		writeError(w, err)
		return
	}

	resp := renderKeys(keys)
	writeSuccess(w, &resp)
}

func (k KeyServer) UpdateKey(w http.ResponseWriter, r *http.Request) {
	req := types.UpdateKeyRequest{}
	err := readRequest(r, &req)
	if err != nil {
		writeError(w, err)
		return
	}

	vars := mux.Vars(r)
	name := vars["name"]
	if name != req.Name {
		writeError(w, errors.New("path and json key names don't match"))
		return
	}

	err = k.manager.Update(req.Name, req.OldPass, req.NewPass)
	if err != nil {
		writeError(w, err)
		return
	}

	key, err := k.manager.Get(req.Name)
	if err != nil {
		writeError(w, err)
		return
	}

	resp := renderKey(key)
	writeSuccess(w, &resp)
}

func (k KeyServer) DeleteKey(w http.ResponseWriter, r *http.Request) {
	req := types.DeleteKeyRequest{}
	err := readRequest(r, &req)
	if err != nil {
		writeError(w, err)
		return
	}

	vars := mux.Vars(r)
	name := vars["name"]
	if name != req.Name {
		writeError(w, errors.New("path and json key names don't match"))
		return
	}

	err = k.manager.Delete(req.Name, req.Passphrase)
	if err != nil {
		writeError(w, err)
		return
	}

	// default code is 0 = success
	resp := types.GenericResponse{
		Log: fmt.Sprintf("Key '%s' deleted", name),
	}
	writeSuccess(w, &resp)
}

func (k KeyServer) Register(r *mux.Router) {
	r.HandleFunc("/", k.GenerateKey).Methods("POST")
	r.HandleFunc("/", k.ListKeys).Methods("GET")
	r.HandleFunc("/{name}", k.GetKey).Methods("GET")
	r.HandleFunc("/{name}", k.UpdateKey).Methods("POST", "PUT")
	r.HandleFunc("/{name}", k.DeleteKey).Methods("DELETE")
}

func renderKey(key lc.KeyInfo) types.KeyResponse {
	return types.KeyResponse{
		Name:    key.Name,
		PubKey:  key.PubKey,
		Address: key.PubKey.Address(),
	}
}

func renderKeys(keys lc.KeyInfos) types.KeyListResponse {
	data := make([]types.KeyResponse, len(keys))
	for i := range keys {
		data[i] = renderKey(keys[i])
	}
	return types.KeyListResponse{
		Keys: data,
	}
}
