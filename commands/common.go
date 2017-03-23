/*
Package commands contains any general setup/helpers valid for all subcommands
*/
package commands

import (
	"fmt"
	"os"
	"path/filepath"

	toml "github.com/pelletier/go-toml" // same as viper, different from tendermint, ugh...
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	keycmd "github.com/tendermint/go-keys/cmd" // these usages can move to some common dir
)

var dirPerm = os.FileMode(0700)

const (
	ChainFlag = "chainid"

	ConfigFile = "config.toml"
	SeedDir    = "seeds"
	ProofDir   = "proofs"
)

// InitCmd will initialize the basecli store
var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the light client for a new chain",
	RunE:  runInit,
}

func init() {
	InitCmd.Flags().Bool("force-reset", false, "DANGEROUS: Wipe clean an existing client store")
}

func AddBasicFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().String(ChainFlag, "", "Chain ID of tendermint node")
}

func runInit(cmd *cobra.Command, args []string) error {
	root := viper.GetString(keycmd.RootFlag)

	if viper.GetBool("force-reset") {
		fmt.Println("Bye, bye data... wiping it all clean!")
		os.RemoveAll(root)
		os.MkdirAll(root, dirPerm)
	}

	err := checkEmpty(root)
	if err != nil {
		return err
	}

	err = initConfigFile(cmd)
	if err != nil {
		return err
	}

	// TODO: accept and validate seed
	return nil
}

func initConfigFile(cmd *cobra.Command) error {
	flags := cmd.Flags()
	tree := toml.TreeFromMap(map[string]interface{}{})

	required := []string{ChainFlag}
	for _, f := range required {
		if !flags.Changed(f) {
			return errors.Errorf(`"--%s" required`, f)
		}
		tree.Set(f, viper.Get(f))
	}

	optional := []string{keycmd.OutputFlag, "encoding"}
	for _, f := range optional {
		if flags.Changed(f) {
			tree.Set(f, viper.Get(f))
		}
	}

	out, err := os.Create(filepath.Join(viper.GetString("root"), ConfigFile))
	if err != nil {
		return errors.WithStack(err)
	}
	defer out.Close()

	// save the config file
	_, err = tree.WriteTo(out)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func checkEmpty(root string) error {
	// we create the keys dir on startup, anything else is a sign of error
	dir, err := os.Open(root)
	if err != nil {
		return errors.WithStack(err)
	}
	defer dir.Close()

	files, err := dir.Readdirnames(-1)
	if err != nil {
		return errors.WithStack(err)
	}

	for _, ours := range []string{ConfigFile, SeedDir, ProofDir} {
		for _, f := range files {
			if f == ours {
				return errors.Errorf(`"%s" already exists, cannot init`, f)
			}
		}
	}
	return nil
}
