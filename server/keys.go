package server

import (
	"net/http"

	"github.com/tendermint/light-client/cryptostore"
)

type KeyServer struct {
	manager cryptostore.Manager
}

func New(manager cryptostore.Manager) *KeyServer {
	return &KeyServer{
		manager: manager,
	}
}

func (k *KeyServer) GenerateKey(w http.ResponseWriter, r *http.Request) {
	req := CreateKeyRequest{}
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

	resp := CreateKeyResponse{}
	writeSuccess(w, &resp)
}

func (k *KeyServer) GenerateSignature(w http.ResponseWriter, r *http.Request) {
	req := GenerateSignatureRequest{}
	err := readRequest(r, &req)
	if err != nil {
		writeError(w, err)
		return
	}

	info, err := k.manager.Get(req.KeyName)
	if err != nil {
		writeError(w, err)
		return
	}

	sig, err := k.manager.Signature(req.KeyName, req.Passphrase, req.Data)
	if err != nil {
		writeError(w, err)
		return
	}

	resp := GenerateSignatureResponse{
		Signature: sig,
		PubKey:    info.PubKey,
		Address:   info.Address,
	}
	writeSuccess(w, &resp)
}

type CreateKeyRequest struct {
	Name       string `json:"name"`
	Passphrase string `json:"password"`
}

type CreateKeyResponse struct{}

type GenerateSignatureRequest struct {
	KeyName    string `json:"name"`
	Passphrase string `json:"password"`
	Data       []byte `json:"data"`
}

type GenerateSignatureResponse struct {
	Signature []byte `json:"signature"`
	PubKey    []byte `json:"pubkey"`
	Address   []byte `json:"address"`
}
