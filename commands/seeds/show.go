package seeds

import (
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/light-client/certifiers"
	"github.com/tendermint/light-client/commands"
)

const (
	heightFlag = "height"
	hashFlag   = "hash"
	fileFlag   = "file"
)

// getCmd represents the get command
var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Show the details of one selected seed",
	Long: `Shows the most recent downloaded key by default.
If desired, you can select by height, validator hash, or a file.
`,
	RunE: showSeed,
}

func init() {
	showCmd.Flags().Int(heightFlag, 0, "Show the seed with closest height to this")
	showCmd.Flags().String(hashFlag, "", "Show the seed matching the validator hash")
	showCmd.Flags().String(fileFlag, "", "Show the seed stored in the given file")
	RootCmd.AddCommand(showCmd)
}

func showSeed(cmd *cobra.Command, args []string) (err error) {
	p := commands.GetProvider()

	var seed certifiers.Seed
	h := viper.GetInt(heightFlag)
	hash := viper.GetString(hashFlag)
	file := viper.GetString(fileFlag)

	// load the seed from the proper place
	if h != 0 {
		seed, err = p.GetByHeight(h)
	} else if hash != "" {
		var vhash []byte
		vhash, err = hex.DecodeString(hash)
		if err == nil {
			seed, err = p.GetByHash(vhash)
		}
	} else if file != "" {
		seed, err = certifiers.LoadSeed(file)
	} else {
		// default is latest seed
		seed, err = certifiers.LatestSeed(p)
	}

	if err != nil {
		return err
	}

	// now render it!
	data, err := json.MarshalIndent(seed, "", "  ")
	fmt.Println(string(data))
	return err
}
