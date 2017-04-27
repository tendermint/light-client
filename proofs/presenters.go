package proofs

import (
	"encoding/hex"

	"github.com/pkg/errors"
	data "github.com/tendermint/go-wire/data"
)

// Presenter allows us to encode queries and parse results in an app-specific way
type Presenter interface {
	MakeKey(string) ([]byte, error)
	ParseData([]byte) (interface{}, error)
}

type Presenters map[string]Presenter

// NewPresenters gives you a default raw presenter
func NewPresenters() Presenters {
	return Presenters{"raw": RawPresenter{}}
}

func (p Presenters) Lookup(app string) (Presenter, error) {
	res, ok := p[app]
	if !ok {
		return nil, errors.Errorf("No presenter registered for %s", app)
	}
	return res, nil
}

var _ Presenter = RawPresenter{}

// RawPresenter just hex-encodes/decodes text.  Useful as default,
// or to embed in other structs for the MakeKey implementation
type RawPresenter struct{}

func (_ RawPresenter) MakeKey(str string) ([]byte, error) {
	return hex.DecodeString(str)
}

func (_ RawPresenter) ParseData(raw []byte) (interface{}, error) {
	return data.Bytes(raw), nil
}
