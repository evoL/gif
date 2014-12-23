package config

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"text/tabwriter"
)

func (c *Config) Print() {
	writer := tabwriter.NewWriter(os.Stdout, 4, 4, 2, ' ', 0)
	defer writer.Flush()

	printStruct(writer, reflect.ValueOf(c).Elem(), 0)
}

func printStruct(writer *tabwriter.Writer, val reflect.Value, level int) {
	t := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)

		for i := 0; i < level; i++ {
			io.WriteString(writer, "\t")
		}

		switch field.Kind() {
		case reflect.String:
			fmt.Fprintf(writer, "%s\t%s\n", t.Field(i).Name, field.String())
		case reflect.Struct:
			fmt.Fprintf(writer, "%s\n", t.Field(i).Name)
			printStruct(writer, field, level+1)
		}
	}
}
