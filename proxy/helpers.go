package proxy

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

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
