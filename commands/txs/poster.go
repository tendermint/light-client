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
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	keycmd "github.com/tendermint/go-crypto/cmd"
	keys "github.com/tendermint/go-crypto/keys"
	lightclient "github.com/tendermint/light-client"
	"github.com/tendermint/light-client/commands"
	"github.com/tendermint/tendermint/rpc/client"
)

type Poster struct {
	name     string
	maker    ReaderMaker
	flagData interface{}
}

type ReaderMaker interface {
	MakeReader() (lightclient.TxReader, error)
	// Flags returns a set of flags to register, as well as a struct
	// which they should parse in to (viper.Unmarshal).  This second
	// argument should be a pointer and will be passed in to TxReader.ReadTxFlags
	Flags() (*flag.FlagSet, interface{})
}

func NewPoster(name string, maker ReaderMaker) *Poster {
	return &Poster{
		name:  name,
		maker: maker,
	}
}

func (p *Poster) CreateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:          p.name,
		RunE:         p.RunE,
		SilenceUsage: true,
	}
	fset, fdata := p.maker.Flags()
	cmd.Flags().AddFlagSet(fset)
	p.flagData = fdata
	return cmd
}

func (p *Poster) RunE(cmd *cobra.Command, args []string) error {
	// get our reader
	reader, err := p.maker.MakeReader()
	if err != nil {
		return err
	}

	// get the pubkey for the tx prep
	name := viper.GetString(NameFlag)
	manager := keycmd.GetKeyManager()
	info, _ := manager.Get(name) // error -> empty pubkey
	pubkey := info.PubKey

	// get input if provided
	input := viper.GetString(InputFlag)
	var tx interface{}
	if input != "" {
		raw, err := readInput(input)
		if err != nil {
			return err
		}

		// parse the input
		tx, err = reader.ReadTxJSON(raw, pubkey)
		if err != nil {
			return err
		}
	} else {
		// we try to parse the flags!
		err := viper.Unmarshal(p.flagData)
		if err != nil {
			return err
		}
		tx, err = reader.ReadTxFlags(p.flagData, pubkey)
		if err != nil {
			return err
		}
	}

	// sign if it is Signable
	var packet []byte
	if sign, ok := tx.(keys.Signable); ok {
		if name == "" {
			return errors.New("--name is required to sign tx")
		}
		packet, err = signTx(manager, sign, name)
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

	js, err := json.MarshalIndent(bres, "", "  ")
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
