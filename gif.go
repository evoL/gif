package main

import (
	"github.com/codegangsta/cli"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "gif"
	app.Usage = "a stupid gif manager"
	app.Action = func(c *cli.Context) {
		HelloWorld()
	}

	app.Run(os.Args)
}
