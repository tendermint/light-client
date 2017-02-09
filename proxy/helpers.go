/*
package proxy provides http handlers to construct a proxy server
for key management, transaction signing, and query validation.

Please read the README and godoc to see how to
configure the server for your application.
*/
package proxy

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	lc "github.com/tendermint/light-client"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

// KeyStore is implemented by cryptostore.Manager
type KeyStore interface {
	lc.KeyManager
	lc.Signer
}

// Node is all rpc stuff together, implemented by rpc.Node
type Node interface {
	lc.Broadcaster
	lc.Checker
	lc.Searcher
}

// RegisterDefault constructs all components and wires them up under
// standard routes.
func RegisterDefault(r *mux.Router, keys KeyStore, node Node,
	txReader lc.SignableReader, valReader lc.ValueReader) {

	key := NewKeyServer(keys)
	sk := r.PathPrefix("/keys").Subrouter()
	key.Register(sk)

	tx := NewTxSigner(node, keys, txReader)
	st := r.PathPrefix("/txs").Subrouter()
	tx.Register(st)

}

func readRequest(r *http.Request, o interface{}) error {
	defer r.Body.Close()
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return errors.Wrap(err, "Read Request")
	}
	err = json.Unmarshal(data, o)
	if err != nil {
		return errors.Wrap(err, "Parse")
	}
	return validate(o)
}

// TODO: much better!!!
func writeError(w http.ResponseWriter, err error) {
	// TODO: better error handling
	w.WriteHeader(500)
	// resp := fmt.Sprintf("%+v", err)
	resp := fmt.Sprintf("%v", err)
	w.Write([]byte(resp))
}

func writeSuccess(w http.ResponseWriter, o interface{}) {
	// TODO: add indent
	data, err := json.Marshal(o)
	if err != nil {
		writeError(w, err)
	} else {
		w.Write(data)
	}
}
