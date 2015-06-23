package image

import (
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"
)

type FlushableWriter interface {
	io.Writer
	Flush() (err error)
}

func DefaultWriter() FlushableWriter {
	return tabwriter.NewWriter(os.Stdout, 4, 4, 2, ' ', 0)
}

func (img *Image) Print() {
	writer := DefaultWriter()
	defer writer.Flush()

	img.PrintTo(writer)
}

func PrintAll(images []Image) {
	writer := DefaultWriter()
	defer writer.Flush()

	for _, img := range images {
		img.PrintTo(writer)
	}
}

func (img *Image) PrintTo(writer io.Writer) {
	fmt.Fprintf(writer, "%s\t", img.Id[:8])

	if img.Url == "" {
		io.WriteString(writer, "local\t")
	} else {
		io.WriteString(writer, "remote\t")
	}

	fmt.Fprintf(writer, "%v\t", img.AddedAt.Format("2006-01-02 15:04:05"))

	if len(img.Tags) > 0 {
		fmt.Fprintln(writer, strings.Join(img.Tags, ", "))
	} else {
		fmt.Fprintln(writer, "(no tags)")
	}
}
