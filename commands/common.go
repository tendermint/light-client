/*
Package commands contains any general setup/helpers valid for all subcommands
*/
package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	toml "github.com/pelletier/go-toml" // same as viper, different from tendermint, ugh...
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	keycmd "github.com/tendermint/go-crypto/cmd" // these usages can move to some common dir
	"github.com/tendermint/light-client/certifiers"
	"github.com/tendermint/light-client/certifiers/client"
	"github.com/tendermint/light-client/certifiers/files"
)

var (
	dirPerm  = os.FileMode(0700)
	provider certifiers.Provider
)

const (
	ChainFlag = "chainid"
	NodeFlag  = "node"
	SeedFlag  = "seed"

	ConfigFile = "config.toml"
)

// InitCmd will initialize the basecli store
var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the light client for a new chain",
	RunE:  runInit,
}

func init() {
	InitCmd.Flags().Bool("force-reset", false, "DANGEROUS: Wipe clean an existing client store")
	InitCmd.Flags().String(SeedFlag, "", "Seed file to import (optional)")
}

func AddBasicFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().String(ChainFlag, "", "Chain ID of tendermint node")
	cmd.PersistentFlags().String(NodeFlag, "", "<host>:<port> to tendermint rpc interface for this chain")
}

func GetProvider() certifiers.Provider {
	if provider == nil {
		// store the keys directory
		rootDir := viper.GetString("root")
		provider = certifiers.NewCacheProvider(
			certifiers.NewMemStoreProvider(),
			files.NewProvider(rootDir),
			client.NewHTTP(viper.GetString(NodeFlag)),
		)
	}
	return provider
}

func GetCertifier() (*certifiers.InquiringCertifier, error) {
	// load up the latest store....
	p := GetProvider()
	// this should get the most recent verified seed
	seed, err := certifiers.LatestSeed(p)
	if err != nil {
		return nil, err
	}
	cert := certifiers.NewInquiring(
		viper.GetString(ChainFlag), seed.Validators, p)
	return cert, nil
}

func runInit(cmd *cobra.Command, args []string) error {
	root := viper.GetString(keycmd.RootFlag)

	if viper.GetBool("force-reset") {
		fmt.Println("Bye, bye data... wiping it all clean!")
		os.RemoveAll(root)
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
	err = initSeed()

	return err
}

func initConfigFile(cmd *cobra.Command) error {
	flags := cmd.Flags()
	tree := toml.TreeFromMap(map[string]interface{}{})

	required := []string{ChainFlag, NodeFlag}
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

func initSeed() (err error) {
	// create a provider....
	p := GetProvider()

	// load a seed file, or get data from the provider
	var seed certifiers.Seed
	seedFile := viper.GetString(SeedFlag)
	if seedFile == "" {
		fmt.Println("Loading validator set from tendermint rpc...")
		seed, err = certifiers.LatestSeed(p)
	} else {
		fmt.Printf("Loading validators from file %s\n", seedFile)
		seed, err = certifiers.LoadSeed(seedFile)
	}
	// can't load the seed? abort!
	if err != nil {
		return err
	}

	// make sure it is a proper seed
	err = seed.ValidateBasic(viper.GetString(ChainFlag))
	if err != nil {
		return err
	}

	// ask the user to verify the validator hash
	fmt.Println("\nImportant: if this is incorrect, all interaction with the chain will be insecure!")
	fmt.Printf("  Given validator hash valid: %X\n", seed.Hash())
	fmt.Println("Is this valid (y/n)?")
	valid := askForConfirmation()
	if !valid {
		return errors.New("Invalid validator hash, try init with proper seed later")
	}

	// if accepted, store seed as current state
	p.StoreSeed(seed)
	return nil
}

func checkEmpty(root string) error {
	// we create the keys dir on startup, anything else is a sign of error
	os.MkdirAll(root, dirPerm)
	dir, err := os.Open(root)
	if err != nil {
		return errors.WithStack(err)
	}
	defer dir.Close()

	files, err := dir.Readdirnames(-1)
	if err != nil {
		return errors.WithStack(err)
	}

	empty := len(files) == 0
	if !empty && len(files) == 1 && files[0] == "keys" {
		empty = true
	}

	if !empty {
		return errors.Errorf(`"%s" contains data, cannot init`, root)
	}
	return nil
}

func askForConfirmation() bool {
	var resp string
	_, err := fmt.Scanln(&resp)
	if err != nil {
		panic(err)
	}
	resp = strings.ToLower(resp)
	if resp == "y" || resp == "yes" {
		return true
	} else if resp == "n" || resp == "no" {
		return false
	} else {
		fmt.Println("Please type yes or no and then press enter:")
		return askForConfirmation()
	}
}
