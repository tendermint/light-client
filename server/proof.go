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

	proof, err := merkle.ReadProof(req.Proof)
	if err != nil {
		writeError(w, errors.Wrap(err, "Loading Proof"))
		return
	}

	resp := ProofResponse{
		RootHash: proof.RootHash,
		Valid:    proof.Verify(req.Key, req.Value, proof.RootHash),
	}
	writeSuccess(w, &resp)
}

type ProofRequest struct {
	Proof []byte `json:"proof"`
	Key   []byte `json:"key"`
	Value []byte `json:"value"`
}

type ProofResponse struct {
	RootHash []byte `json:"roothash"`
	Valid    bool   `json:"valid"`
}
