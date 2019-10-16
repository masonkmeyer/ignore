package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()

	app.Action = func(c *cli.Context) error {
		name := strings.Title(c.Args().Get(0))
		url := fmt.Sprintf("https://raw.githubusercontent.com/github/gitignore/master/%s.gitignore", name)
		resp, err := http.Get(url)

		if err != nil || resp.StatusCode != http.StatusOK {
			return nil
		}

		io.Copy(os.Stdout, resp.Body)
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
