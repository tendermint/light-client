/*
Package commands contains any general setup/helpers valid for all subcommands
*/
package commands

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/tendermint/tmlibs/cli"

	rpcclient "github.com/tendermint/tendermint/rpc/client"

	"github.com/tendermint/light-client/certifiers"
	"github.com/tendermint/light-client/certifiers/client"
	"github.com/tendermint/light-client/certifiers/files"
)

var (
	provider certifiers.Provider
)

const (
	ChainFlag = "chainid"
	NodeFlag  = "node"
)

func AddBasicFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().String(ChainFlag, "", "Chain ID of tendermint node")
	cmd.PersistentFlags().String(NodeFlag, "", "<host>:<port> to tendermint rpc interface for this chain")
}

func GetNode() rpcclient.Client {
	return rpcclient.NewHTTP(viper.GetString(NodeFlag), "/websocket")
}

func GetProvider() certifiers.Provider {
	if provider == nil {
		// store the keys directory
		rootDir := viper.GetString(cli.HomeFlag)
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
