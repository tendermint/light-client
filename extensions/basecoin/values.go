package basecoin

import (
	"fmt"

	bc "github.com/tendermint/basecoin/types"
	wire "github.com/tendermint/go-wire"
	lc "github.com/tendermint/light-client"
	"github.com/tendermint/light-client/tx"
)

type BasecoinValues struct{}

// Turn merkle binary into a json-able struct
func (t BasecoinValues) ReadValue(key, value []byte) (lc.Value, error) {
	// TODO: a more intelligent test that doesn't just brute force
	// all possible types, and allows easier additions to plugins

	// first try to render as an account
	// WTF? We store a pointer, not an object?
	// I had to look at GetAccount() from basecoin/state/state.go to get this to work
	var acct *bc.Account
	err := wire.ReadBinaryBytes(value, &acct)
	if err == nil {
		return renderAccount(acct, value), nil
	}

	// if not render raw
	fmt.Printf("Parse Account: %+v\n", err)
	return tx.NewValue(value), nil
}

func (v BasecoinValues) assertValueReader() lc.ValueReader {
	return v
}

const AccountType = "account"

type Account struct {
	Type  string      `json:"type"`
	Value AccountData `json:"value"` // TODO: custom encoding?
	data  []byte      `json:"-"`
}

type AccountData struct {
	PubKey   tx.JSONPubKey `json:"pub_key,omitempty"` // May be empty, if not known.
	Sequence int           `json:"sequence"`
	Balance  bc.Coins      `json:"coins"`
}

func (a Account) Bytes() []byte {
	return a.data
}

func renderAccount(acct *bc.Account, data []byte) Account {
	return Account{
		Type: AccountType,
		Value: AccountData{
			Sequence: acct.Sequence,
			Balance:  acct.Balance,
			PubKey:   tx.JSONPubKey{acct.PubKey},
		},
		data: data,
	}

}
