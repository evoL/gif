package main

import (
	"bufio"
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/evoL/gif/store"
	"os"
	"regexp"
	"strings"
)

func RecreateCommand(c *cli.Context) {
	really := c.Bool("really")

	if really {
		recreateStore()
	} else {
		fmt.Print("Do you really want to recreate the database? [y/n] ")

		scanner := bufio.NewScanner(os.Stdin)
		if !scanner.Scan() {
			fmt.Println("IO Error: " + scanner.Err().Error())
			os.Exit(1)
		}

		response := strings.TrimSpace(scanner.Text())
		if regexp.MustCompile(`\A[yY]\z`).MatchString(response) {
			recreateStore()
		}
	}
}

func recreateStore() {
	s := getStore()
	defer s.Close()

	migrations := store.DefaultMigrationSource()
	if err := s.Recreate(migrations); err != nil {
		fmt.Println("Error while trying to recreate store: " + err.Error())
		os.Exit(1)
	}
}
