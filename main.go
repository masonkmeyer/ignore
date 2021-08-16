package main

import (
	"github.com/masonkmeyer/ignore/cmd"
	"github.com/masonkmeyer/ignore/repository"
)

func main() {
	provider := repository.NewGitProvider()
	fetch := cmd.NewFetchCmd(provider)
	cmd.Run(fetch)
}
