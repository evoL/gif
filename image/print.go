package image

import (
	"fmt"
	"io"
	"os"
	"text/tabwriter"
)

func (img *Image) Print() {
	writer := tabwriter.NewWriter(os.Stdout, 4, 4, 2, ' ', 0)
	defer writer.Flush()

	img.PrintTo(writer)
}

func (img *Image) PrintTo(writer io.Writer) {
	fmt.Fprintf(writer, "%s\t", img.Id[:8])

	if img.Url == "" {
		io.WriteString(writer, "local\n")
	} else {
		io.WriteString(writer, "remote\n")
	}
}
