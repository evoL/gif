package main

import (
	"bufio"
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/evoL/gif/store"
	"os"
	"strings"
)

func ExportCommand(c *cli.Context) {
	s := getStore()
	defer s.Close()

	var targetFile *os.File
	var err error

	output := c.String("output")
	if output == "-" {
		targetFile = os.Stdout
	} else {
		targetFile, err = os.Create(output)
		if err != nil {
			fmt.Println("Could not create file: " + err.Error())
			os.Exit(1)
		}
	}

	// Detect .zip file extension and enable full export
	exportFiles := c.Bool("zip") || strings.HasSuffix(output, ".zip")

	writer := bufio.NewWriter(targetFile)
	defer writer.Flush()

	err = s.Export(writer, store.NullFilter{}, exportFiles)
	if err != nil {
		fmt.Println("Export error: " + err.Error())
		os.Exit(1)
	}
}
