package server

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
	wire "github.com/tendermint/go-wire"
)

func readRequest(r *http.Request, o interface{}) error {
	defer r.Body.Close()
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return errors.Wrap(err, "Read Request")
	}
	return errors.Wrap(wire.ReadJSONBytes(data, o), "Parse")
}

func writeSuccess(w http.ResponseWriter, o interface{}) {
	data := wire.JSONBytes(o)
	w.Write(data)
}

func writeError(w http.ResponseWriter, err error) {
	// TODO: better error handling
	w.WriteHeader(500)
	// resp := fmt.Sprintf("%+v", err)
	resp := fmt.Sprintf("%v", err)
	w.Write([]byte(resp))
}

func HelloWorld(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>Hello World</h1>")
}
