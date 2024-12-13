package cmd

import (
	"github.com/HappyRedstone/SodiumPlus/builder/internal/all"
	"github.com/spf13/cobra"
)

var BundleAllCommand = &cobra.Command{
	Use:          "bundle",
	Aliases:      []string{"b"},
	Short:        "Create all bundles.",
	Long:         `Create all bundles.`,
	SilenceUsage: true,

	RunE: func(cmd *cobra.Command, args []string) error {
		return all.Bundle()
	},
}

var CleanAllCommand = &cobra.Command{
	Use:          "clean",
	Aliases:      []string{"c"},
	Short:        "Clean all artifacts.",
	Long:         `Clean all artifacts.`,
	SilenceUsage: true,

	RunE: func(cmd *cobra.Command, args []string) error {
		return all.Clean()
	},
}
