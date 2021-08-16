package cmd

import (
	"io"
	"os"
	"strings"

	"github.com/masonkmeyer/ignore/repository"
	"github.com/spf13/cobra"
)

type FetchCmd struct {
	provider repository.Provider
	Find     bool
}

func NewFetchCmd(provider repository.Provider) *FetchCmd {
	return &FetchCmd{
		provider: provider,
	}
}

func (f *FetchCmd) Run() func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		if len(args) == 1 {

			if f.Find {
				result, err := f.provider.Search(args[0])

				if err != nil {
					io.Copy(os.Stderr, strings.NewReader(err.Error()))
					return
				}

				for _, s := range result {
					io.Copy(os.Stdout, strings.NewReader(s))
					io.Copy(os.Stdout, strings.NewReader("\n"))
				}

			} else {

				reader, err := f.provider.Get(args[0])

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
