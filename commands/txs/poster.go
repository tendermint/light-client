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
	keycmd "github.com/tendermint/go-crypto/cmd"
	keys "github.com/tendermint/go-crypto/keys"
	lightclient "github.com/tendermint/light-client"
	"github.com/tendermint/light-client/commands"
	"github.com/tendermint/tendermint/rpc/tendermint/client"
)

type Poster struct {
	name  string
	maker ReaderMaker
	// reader lightclient.TxReader
}

type ReaderMaker interface {
	MakeReader(cmd *cobra.Command, args []string) (lightclient.TxReader, error)
}

func NewPoster(name string, maker ReaderMaker) *Poster {
	return &Poster{
		name:  name,
		maker: maker,
	}
}

func (p *Poster) CreateCommand() *cobra.Command {
	return &cobra.Command{
		Use:  p.name,
		RunE: p.RunE,
	}
}

func (p *Poster) RunE(cmd *cobra.Command, args []string) error {
	fmt.Println("Got", p.name)

	// get our reader
	reader, err := p.maker.MakeReader(cmd, args)
	if err != nil {
		return err
	}

	// get input
	input := viper.GetString(InputFlag)
	if input == "" {
		return errors.New("--input is required")
	}
	raw, err := readInput(input)
	if err != nil {
		return err
	}

	// parse the input
	tx, err := reader.ReadTxJSON(raw)
	if err != nil {
		return err
	}

	// sign if it is Signable
	var packet []byte
	if sign, ok := tx.(keys.Signable); ok {
		name := viper.GetString(NameFlag)
		if name == "" {
			return errors.New("--name is required to sign tx")
		}
		packet, err = signTx(sign, name)
		if err != nil {
			return err
		}
	} else if val, ok := tx.(lightclient.Value); ok {
		packet = val.Bytes()
	} else {
		return errors.Errorf("Reader returned invalid tx type: %#v\n", tx)
	}

	// post the bytes
	endpoint := viper.GetString(commands.NodeFlag)
	node := client.NewHTTP(endpoint, "/websockets")
	bres, err := node.BroadcastTxCommit(packet)
	if err != nil {
		return err
	}

	js, err := json.Marshal(bres)
	if err != nil {
		return err
	}
	fmt.Println(string(js))
	return nil
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
	return ioutil.ReadAll(reader)
}

func signTx(tx keys.Signable, name string) ([]byte, error) {
	manager := keycmd.GetKeyManager()
	prompt := fmt.Sprintf("Please enter passphrase for %s: ", name)
	pass, err := speakeasy.Ask(prompt)
	if err != nil {
		return nil, err
	}
	// TODO: clean up manager so we don't force this
	err = manager.(keys.Signer).Sign(name, pass, tx)
	if err != nil {
		return nil, err
	}
	return tx.TxBytes()
}
