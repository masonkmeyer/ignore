package cmd

import (
	"io"
	"os"
	"strings"

	"github.com/masonkmeyer/ignore/repository"
	"github.com/spf13/cobra"
)

const MAX_RECUSION_DEPTH = 5

type FetchCmd struct {
	provider repository.Provider
	Find     bool
}

func NewFetchCmd(provider repository.Provider) *FetchCmd {
	return &FetchCmd{
		provider: provider,
	}
}

func (f *FetchCmd) Run(aliases map[string]string) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		if len(args) == 1 {

			term := resolveAlias(aliases, args[0], 0)

			if f.Find {
				result, err := f.provider.Search(term)

				if err != nil {
					io.Copy(os.Stderr, strings.NewReader(err.Error()))
					return
				}

				for _, s := range result {
					io.Copy(os.Stdout, strings.NewReader(s))
					io.Copy(os.Stdout, strings.NewReader("\n"))
				}

			} else {
				reader, err := f.provider.Get(term)

				if err != nil {
					io.Copy(os.Stderr, strings.NewReader(err.Error()))
					return
				}

				if reader != nil {
					io.Copy(os.Stdout, reader)
				}
			}
		}
	}
}

func resolveAlias(aliases map[string]string, value string, attempt int) string {
	if attempt > MAX_RECUSION_DEPTH {
		return value
	}

	if val, ok := aliases[value]; ok {
		return resolveAlias(aliases, val, attempt+1)
	}

	return value
}
