package commands

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/tendermint/tmlibs/cli"
	cmn "github.com/tendermint/tmlibs/common"

	"github.com/tendermint/light-client/certifiers"
)

var (
	dirPerm = os.FileMode(0700)
)

const (
	SeedFlag = "seed"
	HashFlag = "valhash"

	ConfigFile = "config.toml"
)

// InitCmd will initialize the basecli store
var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the light client for a new chain",
	RunE:  runInit,
}

var ResetCmd = &cobra.Command{
	Use:   "reset_all",
	Short: "DANGEROUS: Wipe out all client data, including keys",
	RunE:  runResetAll,
}

func init() {
	InitCmd.Flags().Bool("force-reset", false, "Wipe clean an existing client store, except for keys")
	InitCmd.Flags().String(SeedFlag, "", "Seed file to import (optional)")
	InitCmd.Flags().String(HashFlag, "", "Trusted validator hash (must match to accept)")
}

func runInit(cmd *cobra.Command, args []string) error {
	root := viper.GetString(cli.HomeFlag)

	if viper.GetBool("force-reset") {
		resetRoot(root, true)
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

func runResetAll(cmd *cobra.Command, args []string) error {
	root := viper.GetString(cli.HomeFlag)
	resetRoot(root, false)
	return nil
}

func resetRoot(root string, saveKeys bool) {
	tmp := filepath.Join(os.TempDir(), cmn.RandStr(16))
	keys := filepath.Join(root, "keys")
	if saveKeys {
		fmt.Println("Saving private keys", tmp, keys)
		os.Rename(keys, tmp)
	}
	fmt.Println("Bye, bye data... wiping it all clean!")
	os.RemoveAll(root)
	if saveKeys {
		os.Mkdir(root, 0700)
		os.Rename(tmp, keys)
	}
}

type Config struct {
	Chain    string `toml:"chainid,omitempty"`
	Node     string `toml:"node,omitempty"`
	Output   string `toml:"output,omitempty"`
	Encoding string `toml:"encoding,omitempty"`
}

func setConfig(flags *pflag.FlagSet, f string, v *string) {
	if flags.Changed(f) {
		*v = viper.GetString(f)
	}
}

func initConfigFile(cmd *cobra.Command) error {
	flags := cmd.Flags()
	var cfg Config

	required := []string{ChainFlag, NodeFlag}
	for _, f := range required {
		if !flags.Changed(f) {
			return errors.Errorf(`"--%s" required`, f)
		}
	}

	setConfig(flags, ChainFlag, &cfg.Chain)
	setConfig(flags, NodeFlag, &cfg.Node)
	setConfig(flags, cli.OutputFlag, &cfg.Output)
	setConfig(flags, cli.EncodingFlag, &cfg.Encoding)

	out, err := os.Create(filepath.Join(viper.GetString(cli.HomeFlag), ConfigFile))
	if err != nil {
		return errors.WithStack(err)
	}
	defer out.Close()

	// save the config file
	err = toml.NewEncoder(out).Encode(cfg)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func initSeed() (err error) {
	// create a provider....
	trust, source := GetProviders()

	// load a seed file, or get data from the provider
	var seed certifiers.Seed
	seedFile := viper.GetString(SeedFlag)
	if seedFile == "" {
		fmt.Println("Loading validator set from tendermint rpc...")
		seed, err = certifiers.LatestSeed(source)
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

	// validate hash interactively or not
	hash := viper.GetString(HashFlag)
	if hash != "" {
		var hashb []byte
		hashb, err = hex.DecodeString(hash)
		if err == nil && !bytes.Equal(hashb, seed.Hash()) {
			err = errors.Errorf("Seed hash doesn't match expectation: %X", seed.Hash())
		}
	} else {
		err = validateHash(seed)
	}

	if err != nil {
		return err
	}

	// if accepted, store seed as current state
	trust.StoreSeed(seed)
	return nil
}

func validateHash(seed certifiers.Seed) error {
	// ask the user to verify the validator hash
	fmt.Println("\nImportant: if this is incorrect, all interaction with the chain will be insecure!")
	fmt.Printf("  Given validator hash valid: %X\n", seed.Hash())
	fmt.Println("Is this valid (y/n)?")
	valid := askForConfirmation()
	if !valid {
		return errors.New("Invalid validator hash, try init with proper seed later")
	}
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
