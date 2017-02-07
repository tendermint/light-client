package main

import (
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/tendermint/light-client/cryptostore"
	"github.com/tendermint/light-client/cryptostore/filestorage"
	"github.com/tendermint/light-client/server"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

/*
Usage:
  verifier --socket=$SOCKET --keydir=$KEYDIR

Testing:
  curl --unix-socket $SOCKET http://localhost/
  curl -XPOST --unix-socket $SOCKET http://localhost/prove -d '{"proof": "ABCD1234"}'
  curl -XPOST --unix-socket $SOCKET http://localhost/key -d '{"name": "test", "password": "1234"}'
  curl -XPOST --unix-socket $SOCKET http://localhost/sign -d '{"name": "test", "password": "1234", "data": "12345678deadbeef"}'
  curl -XPOST --unix-socket $SOCKET http://localhost/wrap -d '{"name": "test", "password": "1234", "data": "12345678deadbeef"}'
*/

var (
	app    = kingpin.New("verifier", "Local golang server to validate go-merkle proofs")
	socket = app.Flag("socket", "path to unix socket to serve on").String()
	port   = app.Flag("port", "port to listen on").Default("46667").Int()
	keydir = app.Flag("keydir", "directory to store the secret keys").String()
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

	var l net.Listener
	var err error
	if *socket != "" {
		l, err = CreateSocket(*socket)
		app.FatalIfError(err,
			"Cannot create socket: %s", *socket)
	} else {
		l, err = net.Listen("tcp", fmt.Sprintf(":%v", *port))
		app.FatalIfError(err,
			"Cannot listen on port: %v", *port)
	}

	router := mux.NewRouter()
	router.HandleFunc("/", server.HelloWorld).Methods("GET")
	router.HandleFunc("/prove", server.VerifyProof).Methods("POST")

	if keydir != nil && *keydir != "" {
		crypto := cryptostore.New(
			cryptostore.GenEd25519,
			cryptostore.SecretBox,
			filestorage.New(*keydir),
		)
		keystore := server.New(crypto)
		router.HandleFunc("/key", keystore.GenerateKey).Methods("POST")
		router.HandleFunc("/sign", keystore.GenerateSignature).Methods("POST")
	}

	// only set cors for tcp listener
	var h http.Handler
	if *socket == "" {
		allowedHeaders := handlers.AllowedHeaders([]string{"Content-Type"})
		h = handlers.CORS(allowedHeaders)(router)
	} else {
		h = router
	}

	app.FatalIfError(http.Serve(l, h), "Server killed")
}
