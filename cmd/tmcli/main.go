/*
tmcli is an example cli build on light-client that interacts with a
tendermint node running basecoin.

The only basecoin-specific logic should be set up in this main file,
all other packages should support multiple abci apps.

The commands are run in cobra/viper as per tendermint standard

All data is stored in a data dir, set as --data, or default ~/.tmcli

Commands

FIRST: test that current basecoin-proxy command works against v0.9/v0.4 release
 so we can merge to master (also need to rewrite README.md)

* init - takes the chain ID, and and verifies a known seed
  * returns an error if already initialized in that dir, use --root for a new
    dir or --force-reset to wipe existing data clean
* keys - subcommand to run the go-keys cli
* seeds - subcommand to view header/commit/validator seeds
  * show - shows details on one stored seed
  * update - tries to update from known seed to current validator set if possible
    at a minimum it will download current state
  * export - exports some seeds for passing to a peer
  * import - imports seeds from a peer, filling in gaps if the node changed too
      much while you were offline - does lots of validation
    * --dry-run just checks validity but doesn't store (TODO: deeper)
  * TODO: list????
* proofs
  * get - display just as binary or accept plug in to display as json?
  * list????
  * show?
  * export
  * import (--dry-run)

tmcli proof state list --key <k> --height <h>
tmcli proof state get <key> --height <h>
tmcli proof state export <key> --height <h>

tmcli proof tx list -> hashes
tmcli proof tx get <hash>
tmcli proof tx export <hash>

NEXT:
* tx - at least support sending via cli, if not all plugins...
    at this point we would replace the entire basecoin cli

LATER:
* proxy - runs an http server to post and sign tx, make queries, and
    validate merkle proofs.  this will also update seeds in the background
    also show keys and seeds and proofs via HTTP....

--> actually this main program should be in basecoin along with extensions/basecoin
 as an example of easily tools to build a CLI

*/

package main

import (
	"encoding/hex"
	"os"

	"github.com/spf13/cobra"
	btypes "github.com/tendermint/basecoin/types"
	keycmd "github.com/tendermint/go-keys/cmd"
	wire "github.com/tendermint/go-wire"
	"github.com/tendermint/light-client/commands"
	"github.com/tendermint/light-client/commands/proofs"
	"github.com/tendermint/light-client/commands/seeds"
)

// BaseCli represents the base command when called without any subcommands
var BaseCli = &cobra.Command{
	Use:   "tmcli",
	Short: "Light client for tendermint",
	Long: `Tmcli is a full-fledged, but generic light-client app for tendermint.

You can manager keys, sync validator sets, requests proofs, and
post transactions. All functionality exposed as a cli tool as well as
over a JSON API.

This works with raw hex-encoded bytes for transactions and state data.
It is intended to be imported in a specific abci app and customized with
some parsing code to enable a customized cli that is aware of the
app-specific data structures.
`,
}

// TODO: remove this. basecoin account viewer for demo
type AccountPresenter struct{}

func (_ AccountPresenter) MakeKey(str string) ([]byte, error) {
	res, err := hex.DecodeString(str)
	if err == nil {
		res = append([]byte("base/a/"), res...)
	}
	return res, err
}

func (_ AccountPresenter) ParseData(raw []byte) (interface{}, error) {
	var acc *btypes.Account
	err := wire.ReadBinaryBytes(raw, &acc)
	return acc, err
}

func init() {
	commands.AddBasicFlags(BaseCli)

	// set up the various commands to use
	BaseCli.AddCommand(keycmd.RootCmd)
	BaseCli.AddCommand(commands.InitCmd)
	BaseCli.AddCommand(seeds.RootCmd)
	// TODO: when subclassing register some parsers with
	// proofs.StatePresenters["app"] = pres
	// proofs.TxPresenters["app"] = pres
	proofs.StatePresenters["account"] = AccountPresenter{}
	BaseCli.AddCommand(proofs.RootCmd)
}

func main() {
	keycmd.PrepareMainCmd(BaseCli, "TM", os.ExpandEnv("$HOME/.tmcli"))
	BaseCli.Execute()
	// err := BaseCli.Execute()
	// if err != nil {
	// 	fmt.Printf("%+v\n", err)
	// }
}
