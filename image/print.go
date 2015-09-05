package image

import (
	"fmt"
	"github.com/dustin/go-humanize"
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
	return WriterFor(os.Stdout)
}

func WriterFor(writer io.Writer) FlushableWriter {
	return tabwriter.NewWriter(writer, 4, 4, 2, ' ', 0)
}

func (img *Image) Print() {
	writer := DefaultWriter()
	defer writer.Flush()

	img.PrintTo(writer, false)
}

func PrintAll(images []Image) {
	writer := DefaultWriter()
	defer writer.Flush()

	PrintAllTo(images, writer)
}

func PrintAllTo(images []Image, writer io.Writer) {
	for _, img := range images {
		img.PrintTo(writer, false)
	}
}

func (img *Image) PrintTo(writer io.Writer, flush bool) {
	fmt.Fprintf(writer, "%s\t", img.Id[:8])

	if img.Url == "" {
		io.WriteString(writer, "local\t")
	} else {
		io.WriteString(writer, "remote\t")
	}

	fmt.Fprintf(writer, "%s\t", humanize.Bytes(img.Size))

	fmt.Fprintf(writer, "%v\t", img.AddedAt.Format("2006-01-02 15:04:05"))

	if len(img.Tags) > 0 {
		fmt.Fprint(writer, strings.Join(img.Tags, ", "))
	} else {
		fmt.Fprint(writer, "(no tags)")
	}

	if flush {
		io.WriteString(writer, "\f")
	} else {
		io.WriteString(writer, "\n")
	}
}
