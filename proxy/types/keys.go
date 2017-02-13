package types

import "github.com/tendermint/light-client/tx"

// CreateKeyRequest is sent to create a new key
type CreateKeyRequest struct {
	Name       string `json:"name" validate:"required,min=4,printascii"`
	Passphrase string `json:"passphrase" validate:"required,min=10"`
}

type DeleteKeyRequest CreateKeyRequest

// UpdateKeyRequest is sent to update the passphrase for an existing key
type UpdateKeyRequest struct {
	Name    string `json:"name" validate:"required,min=4,printascii"`
	OldPass string `json:"passphrase"  validate:"required,min=10"`
	NewPass string `json:"new_passphrase" validate:"required,min=10"`
}

// KeyResponse returns public info on one key
type KeyResponse struct {
	Name    string     `json:"name"`
	PubKey  tx.HexData `json:"pub_key"` // TODO: return in [byte, string] format?
	Address tx.HexData `json:"address"`
}

// KeyListResponse returns info on all keys in the store
type KeyListResponse struct {
	Keys []KeyResponse `json:"keys"`
}
