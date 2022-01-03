package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/masonkmeyer/ignore/cmd"
	"github.com/masonkmeyer/ignore/repository"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigName(".ignore")
	viper.AddConfigPath("$HOME")
	viper.SetConfigType("yaml")
	viper.SetDefault("aliases", make(map[string]string))

	if err := getOrCreateConfig(); err != nil {
		io.Copy(os.Stderr, strings.NewReader(err.Error()))
		return
	}

	provider := repository.NewGitProvider()
	fetch := cmd.NewFetchCmd(provider)
	alias := cmd.NewAliasCmd()
	cmd.Run(fetch, alias)
}

func getOrCreateConfig() error {
	if err := viper.ReadInConfig(); err != nil {
		if err, ok := err.(viper.ConfigFileNotFoundError); ok {
			dirname, err := os.UserHomeDir()

			if err != nil {
				return err
			}

			if _, err := os.Create(fmt.Sprintf("%s/.ignore", dirname)); err != nil {
				return err
			}
		} else {
			return err
		}
	}

	return nil
}
