package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Run(fetchCmd *FetchCmd, aliasCmd *AliasCmd) error {
	aliases := viper.GetStringMapString(ALIASES)

	var rootCmd = &cobra.Command{
		Use:              "ignore [search term]",
		Short:            "ignore is a tool to help find fetch .gitignore files",
		Run:              fetchCmd.Run(aliases),
		Args:             cobra.MinimumNArgs(1),
		TraverseChildren: true,
	}

	rootCmd.Flags().BoolVarP(&fetchCmd.Find, "find", "f", false, "Search for available gitignore files")

	var alias = &cobra.Command{
		Use:     "alias",
		Short:   "Alias a .gitignore file name",
		Example: "ignore alias mac Global/macOS",
		Args:    cobra.MinimumNArgs(2),
		Run:     aliasCmd.Alias(aliases),
	}

	var unalias = &cobra.Command{
		Use:     "unalias",
		Short:   "Unalias a .gitignore file name",
		Example: "ignore alias mac",
		Args:    cobra.MinimumNArgs(1),
		Run:     aliasCmd.Unalias(aliases),
	}

	rootCmd.AddCommand(alias)
	rootCmd.AddCommand(unalias)
	return rootCmd.Execute()
}
