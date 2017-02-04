package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"os"
	"text/tabwriter"
)

func main() {
	w := new(tabwriter.Writer)

	// Format in tab-separated columns with a tab stop of 8.
	w.Init(os.Stdout, 1, 8, 1, '\t', 0)
	fmt.Fprintf(w, "file\timprov\torig siz\tcomp siz\n")
	for i := 1; i < 256; i++ {
		fname := fmt.Sprintf("%d.json", i)
		b, err := ioutil.ReadFile(fmt.Sprintf("fixtures/%s", fname))
		if err != nil {
			panic(err)
		}
		var buf bytes.Buffer
		c := gzip.NewWriter(&buf)
		c.Write(b)
		c.Close()
		savings := (float64(len(b)) - float64(buf.Len())) / float64(len(b)) * 100
		fmt.Fprintf(w, "%s\t%f\t%d bytes\t%d bytes\n", fname, savings, len(b), buf.Len())
	}
	w.Flush()
}
