/*
package proxy provides http handlers to construct a proxy server
for key management, transaction signing, and query validation.

Please read the README and godoc to see how to
configure the server for your application.
*/
package proxy

import (
	"fmt"
	"io/ioutil"
	"net/http"

	wire "github.com/tendermint/go-wire"
	lc "github.com/tendermint/light-client"
	"github.com/tendermint/light-client/proxy/types"
	"github.com/tendermint/light-client/rpc"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

// KeyStore is implemented by cryptostore.Manager
type KeyStore interface {
	lc.KeyManager
	lc.Signer
}

// RegisterDefault constructs all components and wires them up under
// standard routes.
//
// TODO: something more intelligent for getting validators,
// this is pretty insecure right now
func RegisterDefault(r *mux.Router, keys KeyStore, node rpc.Node,
	txReader lc.SignableReader, valReader lc.ValueReader) {

	key := NewKeyServer(keys)
	sk := r.PathPrefix("/keys").Subrouter()
	key.Register(sk)

	tx := NewTxSigner(node, keys, txReader)
	st := r.PathPrefix("/txs").Subrouter()
	tx.Register(st)

	// query the node for the validator - soon at least cache locally
	vals, err := node.Validators()
	if err != nil {
		panic(err)
	}
	cert := rpc.StaticCertifier{Vals: vals.Validators}

	view := NewViewer(node, node, cert)
	view.Register(r)
}

func readRequest(r *http.Request, o interface{}) error {
	defer r.Body.Close()
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return errors.Wrap(err, "Read Request")
	}
	err = wire.ReadJSONBytes(data, o)
	// err = json.Unmarshal(data, o)
	if err != nil {
		return errors.Wrap(err, "Parse")
	}
	return validate(o)
}

// most errors are bad input, so 406... do better....
func writeError(w http.ResponseWriter, err error) {
	fmt.Printf("\nError: %+v\n\n", err)
	res := types.GenericResponse{
		Code: 406,
		// Log:  fmt.Sprintf("%+v", err),
		Log: err.Error(),
	}
	writeCode(w, &res, 406)
}

func writeCode(w http.ResponseWriter, o interface{}, code int) {
	// data, err := json.MarshalIndent(o, "", "    ")
	// if err != nil {
	// 	writeError(w, err)
	// 	return
	// }
	data := wire.JSONBytesPretty(o)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}

func writeSuccess(w http.ResponseWriter, o interface{}) {
	writeCode(w, o, 200)
}
