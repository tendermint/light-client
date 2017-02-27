package basecoin

import (
	"bytes"

	"github.com/pkg/errors"
	bc "github.com/tendermint/basecoin/types"
	wire "github.com/tendermint/go-wire"
	lc "github.com/tendermint/light-client"
	"github.com/tendermint/light-client/tx"
)

type BasecoinValues struct {
	readers []lc.ValueReader
}

func NewBasecoinValues() *BasecoinValues {
	val := BasecoinValues{}
	val.RegisterPlugin(accountParser{})
	return &val
}

// Turn merkle binary into a json-able struct
func (t *BasecoinValues) ReadValue(key, value []byte) (lc.Value, error) {
	// try all plugins in order until a match
	// they should check the key to see if they are appropriate
	for _, read := range t.readers {
		val, err := read.ReadValue(key, value)
		if err == nil {
			return val, err
		}
	}

	// if not render raw
	return tx.NewValue(value), nil
}

func (t *BasecoinValues) RegisterPlugin(reader lc.ValueReader) {
	t.readers = append(t.readers, reader)
}

type accountParser struct{}

func (_ accountParser) ReadValue(key, value []byte) (lc.Value, error) {
	if len(key) == 0 || bytes.Equal([]byte("base/a/"), key[:7]) {
		// first try to render as an account
		// WTF? We store a pointer, not an object?
		// I had to look at GetAccount() from basecoin/state/state.go to get this to work
		var acct *bc.Account
		err := wire.ReadBinaryBytes(value, &acct)
		if err == nil {
			return renderAccount(acct, value), nil
		}
		return nil, errors.Wrap(err, "Parsing account")
	}
	return nil, errors.New("Ignoring this key")
}

func (v *BasecoinValues) assertValueReader() lc.ValueReader {
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
