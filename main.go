package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/urfave/cli"
)

const suffix = ".gitignore"

// get will fetch the ignore from github's gitignore repository.
// It will attempt to get the package by name. It will follow 3 symlinks.
func get(lang string, attempts uint8) string {
	url := fmt.Sprintf("https://raw.githubusercontent.com/github/gitignore/master/%s.gitignore", lang)
	resp, err := http.Get(url)

	if err != nil || resp.StatusCode != http.StatusOK {
		return ""
	}

	b, _ := ioutil.ReadAll(resp.Body)

	result := string(b)

	// If we only have a single line nad the line ends with .gitignore
	lines := strings.Count(result, "\n")

	if strings.Contains(result, suffix) && lines == 0 && attempts <= 3 {
		name := strings.Replace(result, suffix, "", 1)
		attempts++
		return get(name, attempts)
	}

	return result
}

func main() {
	app := cli.NewApp()

	app.Action = func(c *cli.Context) error {
		name := strings.Title(c.Args().Get(0))
		resp := strings.NewReader(get(name, 1))

		io.Copy(os.Stdout, resp)
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
