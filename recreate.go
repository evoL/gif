package main

import (
	"bufio"
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/evoL/gif/store"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

func RecreateCommand(c *cli.Context) {
	really := c.Bool("really")
	verbose := c.Bool("verbose")

	if really {
		recreateStore(verbose)
	} else {
		fmt.Print("Do you really want to recreate the database? [y/n] ")

		scanner := bufio.NewScanner(os.Stdin)
		if !scanner.Scan() {
			fmt.Println("IO Error: " + scanner.Err().Error())
			os.Exit(1)
		}

		response := strings.TrimSpace(scanner.Text())
		if regexp.MustCompile(`\A[yY]\z`).MatchString(response) {
			recreateStore(verbose)
		}
	}
}

func recreateStore(verbose bool) {
	s := getStore()
	defer s.Close()

	// Step 1: Make a backup

	if verbose {
		fmt.Println("Exporting backup…")
	}

	backupFile, err := ioutil.TempFile("", "gif-backup")
	if err != nil {
		fmt.Println("Could not create file: " + err.Error())
		os.Exit(1)
	}

	writer := bufio.NewWriter(backupFile)

	err = s.Export(writer, store.NullFilter{}, false)
	if err != nil {
		fmt.Println("Export error: " + err.Error())
		os.Exit(1)
	}

	writer.Flush()

	if verbose {
		fmt.Println("Backup written to " + backupFile.Name())
	}

	// Step 2: Drop the schema

	if verbose {
		fmt.Println("Dropping the schema…")
	}

	if err = s.Implode(); err != nil {
		fmt.Printf("Error when trying to drop the schema: %v\nBackup is available at %v\n", err.Error(), backupFile.Name())
		os.Exit(1)
	}

	// Step 3: Create the schema anew

	if verbose {
		fmt.Println("Recreating the schema…")
	}

	if err = s.Initialize(); err != nil {
		fmt.Printf("Error when trying to initialize the schema: %v\nBackup is available at %v\n", err.Error(), backupFile.Name())
		os.Exit(1)
	}

	// Step 4: Import!

	if verbose {
		fmt.Println("Importing the files…")
	}
}
