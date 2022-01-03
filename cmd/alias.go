package cmd

import (
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const ALIASES = "aliases"

type AliasCmd struct {
}

func NewAliasCmd() *AliasCmd {
	return &AliasCmd{}
}

func (a *AliasCmd) Alias(aliases map[string]string) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		aliases[args[0]] = args[1]
		a.save(aliases)
	}
}

func (a *AliasCmd) Unalias(aliases map[string]string) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		delete(aliases, args[0])
		a.save(aliases)
	}
}

func (a *AliasCmd) save(aliases map[string]string) {
	viper.Set(ALIASES, aliases)
	err := viper.WriteConfig()
	if err != nil {
		io.Copy(os.Stderr, strings.NewReader(err.Error()))
	}
}
