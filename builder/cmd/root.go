package cmd

import (
	"os"

	"github.com/NoSadBeHappy/SodiumPlus/builder/cmd/client"
	"github.com/NoSadBeHappy/SodiumPlus/builder/cmd/server"
	"github.com/NoSadBeHappy/SodiumPlus/builder/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:              "builder",
	Short:            "A build system for Packwiz.",
	Long:             `A build system for Packwiz.`,
	TraverseChildren: true,
}

func Execute() {
	err := rootCmd.Execute()

	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringP("pack-file", "f", "pack.toml", "The pack file to use")
	rootCmd.Flags().StringP("server-dir", "d", "server", "The server template folder")
	rootCmd.Flags().StringP("config-file", "c", "build.toml", "The path to the config file")

	viper.BindPFlag("pack-file", rootCmd.Flags().Lookup("pack-file"))
	viper.SetDefault("pack-file", "pack.toml")

	viper.BindPFlag("server-folder", rootCmd.Flags().Lookup("server-dir"))
	viper.SetDefault("server-folder", "server")

	viper.BindPFlag("config-file", rootCmd.Flags().Lookup("config-file"))
	viper.SetDefault("config-file", "build.toml")

	rootCmd.AddCommand(client.ClientCommand)

	cfg, err := config.GetConfig()

	if err != nil {
		panic(err)
	}

	if cfg.Server.Enabled {
		rootCmd.AddCommand(server.ServerCommand)
	}

	rootCmd.AddCommand(ListCommand)
	rootCmd.AddCommand(CleanAllCommand)
	rootCmd.AddCommand(BundleAllCommand)
	rootCmd.AddCommand(ServeCommand)
}
