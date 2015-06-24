package main

import (
	"bufio"
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/evoL/gif/store"
	"os"
)

func ExportCommand(c *cli.Context) {
	s := getStore()
	defer s.Close()

	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	err := s.Export(writer, store.NullFilter{}, false)
	if err != nil {
		fmt.Println("Export error: " + err.Error())
		os.Exit(1)
	}
}
