package cmd

import (
	"github.com/spf13/cobra"
)

func Run(fetchCmd *FetchCmd) error {
	var rootCmd = &cobra.Command{
		Use:              "ignore [search term]",
		Short:            "ignore is a tool to help find fetch .gitignore files",
		Run:              fetchCmd.Run(),
		Args:             cobra.MinimumNArgs(1),
		TraverseChildren: true,
	}

	rootCmd.Flags().BoolVarP(&fetchCmd.Find, "find", "f", false, "Search for available gitignore files")

	return rootCmd.Execute()
}
