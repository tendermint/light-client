/*
basecli is an example cli build on light-client that interacts with a
tendermint node running basecoin.

The only basecoin-specific logic should be set up in this main file,
all other packages should support multiple abci apps.

The commands are run in cobra/viper as per tendermint standard

All data is stored in a data dir, set as --data, or default ~/.basecli

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
  * validate - verifies and shows details of file (just import --dry-run?)
* proofs - like seeds, store them (later)
  * list
  * show
  * export
  * import
  * ??get a new proof?? this makes a query via command line

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
	"os"

	"github.com/spf13/cobra"
	keycmd "github.com/tendermint/go-keys/cmd"
	"github.com/tendermint/light-client/commands"
	"github.com/tendermint/light-client/commands/seeds"
)

// BaseCli represents the base command when called without any subcommands
var BaseCli = &cobra.Command{
	Use:   "basecli",
	Short: "Light client for basecoin",
	Long: `Basecli is a full-fledged light-client app for basecoin.

You can manager keys, sync validator sets, requests proofs, and
post transactions. All functionality exposed as a cli tool as well as
over a JSON API.`,
}

func init() {
	commands.AddBasicFlags(BaseCli)

	// set up the various commands to use
	BaseCli.AddCommand(keycmd.RootCmd)
	BaseCli.AddCommand(commands.InitCmd)
	BaseCli.AddCommand(seeds.RootCmd)
}

func main() {
	keycmd.PrepareMainCmd(BaseCli, "TM", os.ExpandEnv("$HOME/.basecli"))
	BaseCli.Execute()
	// err := BaseCli.Execute()
	// if err != nil {
	// 	fmt.Printf("%+v\n", err)
	// }
}
