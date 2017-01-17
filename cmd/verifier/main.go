package main

import (
	"net"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/tendermint/light-client/server"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

/*
Usage:
  verifier --socket=<socket>

Testing:
  curl --unix-socket <socket> http://localhost/
  curl -XPOST --unix-socket <socket> http://localhost/prove -d '{"proof": "ABCD1234"}'
*/

var (
	app    = kingpin.New("verifier", "Local golang server to validate go-merkle proofs")
	socket = app.Flag("socket", "path to unix socket to server on").Required().String()
)

// CreateSocket deletes existing socket if there, creates a new one,
// starts a server on the socket, and sets permissions to 0600
func CreateSocket(socket string) (net.Listener, error) {
	err := os.Remove(socket)
	if err != nil && !os.IsNotExist(err) {
		// only fail if it does exist and cannot be deleted
		return nil, err
	}

	l, err := net.Listen("unix", socket)
	if err != nil {
		return nil, err
	}

	mode := os.FileMode(0700) | os.ModeSocket
	err = os.Chmod(socket, mode)
	if err != nil {
		l.Close()
		return nil, err
	}

	return l, nil
}

func main() {
	kingpin.MustParse(app.Parse(os.Args[1:]))
	l, err := CreateSocket(*socket)
	app.FatalIfError(err,
		"Cannot create socket: %s", *socket)

	router := mux.NewRouter()
	router.HandleFunc("/", server.HelloWorld).Methods("GET")
	router.HandleFunc("/prove", server.VerifyProof).Methods("POST")

	app.FatalIfError(http.Serve(l, router), "Server killed")
}
