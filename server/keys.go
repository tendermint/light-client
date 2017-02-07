package server

import (
	"net/http"

	"github.com/tendermint/light-client/cryptostore"
	"github.com/tendermint/light-client/mock"
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

	tx := &mock.OneSign{Data: req.Data}
	err = k.manager.Sign(req.KeyName, req.Passphrase, tx)
	if err != nil {
		writeError(w, err)
		return
	}

	pk := tx.PubKey
	resp := GenerateSignatureResponse{
		Signature: tx.Sig.Bytes(),
		PubKey:    pk.Bytes(),
		Address:   pk.Address(),
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
