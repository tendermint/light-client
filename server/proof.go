package server

import (
	"net/http"

	"github.com/pkg/errors"
	merkle "github.com/tendermint/go-merkle"
)

// VerifyProof is a handler to embed in your local server
func VerifyProof(w http.ResponseWriter, r *http.Request) {
	req := ProofRequest{}
	err := readRequest(r, &req)
	if err != nil {
		writeError(w, err)
		return
	}

	proof, err := merkle.LoadProof(req.Proof)
	if err != nil {
		writeError(w, errors.Wrap(err, "Loading Proof"))
		return
	}

	resp := ProofResponse{
		Key:      proof.Key(),
		Value:    proof.Value(),
		RootHash: proof.Root(),
		Valid:    proof.Valid(),
	}
	writeSuccess(w, &resp)
}

type ProofRequest struct {
	Proof []byte `json:"proof"`
}

type ProofResponse struct {
	Key      []byte `json:"key"`
	Value    []byte `json:"value"`
	RootHash []byte `json:"roothash"`
	Valid    bool   `json:"valid"`
}
