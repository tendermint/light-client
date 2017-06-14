package txs

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/bgentry/speakeasy"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	crypto "github.com/tendermint/go-crypto"
	keycmd "github.com/tendermint/go-crypto/cmd"
	"github.com/tendermint/go-crypto/keys"

	ctypes "github.com/tendermint/tendermint/rpc/core/types"

	lightclient "github.com/tendermint/light-client"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "tx",
	Short: "Create and post transactions to the node",
}

const (
	NameFlag  = "name"
	InputFlag = "input"
)

func init() {
	RootCmd.PersistentFlags().String(NameFlag, "", "name to sign the tx")
	RootCmd.PersistentFlags().String(InputFlag, "", "file with tx in json format")
}

// GetSigner returns the pub key that will sign the tx
// returns empty key if no name provided
func GetSigner() crypto.PubKey {
	name := viper.GetString(NameFlag)
	manager := keycmd.GetKeyManager()
	info, _ := manager.Get(name) // error -> empty pubkey
	return info.PubKey
}

// Sign if it is Signable, otherwise, just convert it to bytes
func Sign(tx interface{}) (packet []byte, err error) {
	name := viper.GetString(NameFlag)
	manager := keycmd.GetKeyManager()

	if sign, ok := tx.(keys.Signable); ok {
		if name == "" {
			return nil, errors.New("--name is required to sign tx")
		}
		packet, err = signTx(manager, sign, name)
	} else if val, ok := tx.(lightclient.Value); ok {
		packet = val.Bytes()
	} else {
		err = errors.Errorf("Reader returned invalid tx type: %#v\n", tx)
	}
	return
}

// LoadJSON will read a json file from disk if --input is passed in
// template is a pointer to a struct that can hold the expected data (&MyTx{})
//
// If not data is provided, returns (nil, nil)
// If data is provided and passes, returns (template, nil)
// If data is provided but not parsable, returns (nil, err)
func LoadJSON(template interface{}) (interface{}, error) {
	input := viper.GetString(InputFlag)
	if input == "" {
		return nil, nil
	}

	// load the input
	raw, err := readInput(input)
	if err != nil {
		return nil, err
	}

	// parse the input
	err = json.Unmarshal(raw, template)
	if err != nil {
		return nil, err
	}
	return template, nil
}

// OutputTx prints the tx result to stdout
// TODO: something other than raw json?
func OutputTx(res *ctypes.ResultBroadcastTxCommit) error {
	js, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(js))
	return nil
}

func signTx(manager keys.Manager, tx keys.Signable, name string) ([]byte, error) {
	prompt := fmt.Sprintf("Please enter passphrase for %s: ", name)
	pass, err := speakeasy.Ask(prompt)
	if err != nil {
		return nil, err
	}
	err = manager.Sign(name, pass, tx)
	if err != nil {
		return nil, err
	}
	return tx.TxBytes()
}

func readInput(file string) ([]byte, error) {
	var reader io.Reader
	// get the input stream
	if file == "-" {
		reader = os.Stdin
	} else {
		f, err := os.Open(file)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		reader = f
	}

	// and read it all!
	data, err := ioutil.ReadAll(reader)
	return data, errors.WithStack(err)
}
