package basecoin

import (
	"github.com/pkg/errors"
	bc "github.com/tendermint/basecoin/types"
	crypto "github.com/tendermint/go-crypto"
	keys "github.com/tendermint/go-keys"
	wire "github.com/tendermint/go-wire"
)

type AppTx struct {
	chainID string
	signer  crypto.PubKey
	Tx      *bc.AppTx
}

func (a *AppTx) assertSignable() keys.Signable {
	return a
}

func (t BasecoinTx) readAppTx(data []byte) (keys.Signable, error) {
	tx, err := parseAppTx(data, t.appData)
	app := AppTx{
		chainID: t.chainID,
		Tx:      tx,
	}
	return &app, err
}

// SignBytes returned the unsigned bytes, needing a signature
func (a *AppTx) SignBytes() []byte {
	return a.Tx.SignBytes(a.chainID)
}

// Sign will add a signature and pubkey.
//
// Depending on the Signable, one may be able to call this multiple times for multisig
// Returns error if called with invalid data or too many times
func (a *AppTx) Sign(pubkey crypto.PubKey, sig crypto.Signature) error {
	if a.signer != nil {
		return errors.New("AppTx already signed")
	}
	a.Tx.SetSignature(sig)
	a.signer = pubkey
	return nil
}

// Signers will return the public key(s) that signed if the signature
// is valid, or an error if there is any issue with the signature,
// including if there are no signatures
func (a *AppTx) Signers() ([]crypto.PubKey, error) {
	if a.signer == nil {
		return nil, errors.New("No signatures on AppTx")
	}
	return []crypto.PubKey{a.signer}, nil
}

// TxBytes returns the transaction data as well as all signatures
// It should return an error if Sign was never called
func (a *AppTx) TxBytes() ([]byte, error) {
	// TODO: verify it is signed

	// Code and comment from: basecoin/cmd/commands/tx.go
	// Don't you hate having to do this?
	// How many times have I lost an hour over this trick?!
	txBytes := wire.BinaryBytes(struct {
		bc.Tx `json:"unwrap"`
	}{a.Tx})
	return txBytes, nil
}
