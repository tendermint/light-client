/*
package types contains all json structs used by the proxy server for input
and output.
*/
package types

import (
	"encoding/json"
)

// import "gopkg.in/go-playground/validator.v9"
// err := validate.Struct(mystruct)
// validationErrors := err.(validator.ValidationErrors)

// PostTxRequest is sent to sign and post a new transaction
type PostTxRequest struct {
	Name       string          `json:"name" validate:"required,min=4"`
	Passphrase string          `json:"passphrase" validate:"required,min=10"`
	Data       json.RawMessage `json:"data" validate:"required"` // this is handled by SignableReader
}

// QueryTxResponse is returned on success (GenericResponse on failure)
// Also returned for proofs, with Proven = true
type QueryTxResponse struct {
	Height uint64          `json:"height"`
	Key    []byte          `json:"key"`    // TODO: make sure this is hex encoded
	Value  json.RawMessage `json:"value"`  // this is from ValueReader
	Proven bool            `json:"proven"` // only true if we verified all headers
}
