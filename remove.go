package main

import (
	"bufio"
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/evoL/gif/image"
	"github.com/evoL/gif/store"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func RemoveCommand(c *cli.Context) {
	all := c.Bool("all")
	really := c.Bool("really")

	if !c.Args().Present() && !all {
		cli.ShowCommandHelp(c, "remove")
		os.Exit(1)
	}

	s := getStore()
	defer s.Close()

	filter := listFilter(c)

	images, err := s.List(filter)
	if err != nil {
		fmt.Println("Error while fetching: " + err.Error())
		os.Exit(1)
	}

	imageCount := len(images)

	switch {
	// case imageCount == 0: do nothing
	case imageCount == 1:
		removeSingle(s, &images[0], really)
	case imageCount > 1:
		if all {
			removeAll(s, images, really)
		} else {
			removeMultiple(s, images)
		}
	}
}

func removeSingle(s *store.Store, img *image.Image, really bool) {
	if really {
		if err := s.Remove(img); err != nil {
			fmt.Println("Error while removing image: " + err.Error())
			os.Exit(1)
		}
	} else {
		img.Print()
		fmt.Println()
		fmt.Print("Do you really want to remove this image? [y/n] ")

		scanner := bufio.NewScanner(os.Stdin)
		if !scanner.Scan() {
			fmt.Println("IO Error: " + scanner.Err().Error())
			os.Exit(1)
		}

		response := strings.TrimSpace(scanner.Text())
		if regexp.MustCompile(`\A[yY]\z`).MatchString(response) {
			if err := s.Remove(img); err != nil {
				fmt.Println("Error while removing image: " + err.Error())
				os.Exit(1)
			}
		}
	}
}

func removeAll(s *store.Store, images []image.Image, really bool) {
	if really {
		if err := s.RemoveAll(images); err != nil {
			fmt.Println("Error while removing image: " + err.Error())
			os.Exit(1)
		}
	} else {
		fmt.Printf("%v images\n", len(images))

		image.PrintAll(images)

		fmt.Println()
		fmt.Print("Do you really want to remove those images? [y/n] ")

		scanner := bufio.NewScanner(os.Stdin)
		if !scanner.Scan() {
			fmt.Println("IO Error: " + scanner.Err().Error())
			os.Exit(1)
		}

		response := strings.TrimSpace(scanner.Text())
		if regexp.MustCompile(`\A[yY]\z`).MatchString(response) {
			if err := s.RemoveAll(images); err != nil {
				fmt.Println("Error while removing image: " + err.Error())
				os.Exit(1)
			}
		}
	}
}

func removeMultiple(s *store.Store, images []image.Image) {
	fmt.Printf("%v images\n", len(images))

	writer := image.DefaultWriter()

	for i, img := range images {
		io.WriteString(writer, fmt.Sprintf("%v\t", i+1))
		img.PrintTo(writer, false)
	}
	writer.Flush()
	fmt.Println()

	fmt.Println("Select images you want to remove: (Enter comma-separated numbers)")
	fmt.Print("=> ")

	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		fmt.Println()
		os.Exit(1)
	}

	// Validate input
	numberList := strings.TrimSpace(scanner.Text())
	if !regexp.MustCompile(`\A\d+(\s*,\s*\d+)*\z`).MatchString(numberList) {
		fmt.Println("Invalid input.")
		os.Exit(1)
	}

	// Split into actual slice of strings
	numbers := regexp.MustCompile(`\s*,\s*`).Split(numberList, -1)

	for _, stringNumber := range numbers {
		number, _ := strconv.Atoi(stringNumber)
		index := number - 1

		if err := s.Remove(&images[index]); err != nil {
			fmt.Println("Error while removing image: " + err.Error())
		}
	}
}
