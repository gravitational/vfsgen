package vfsgen

import (
	"fmt"
	"io"
)

// byteWriter encodes the input as a []byte literal.
// It tracks the total number of bytes written.
type byteWriter struct {
	w io.Writer
	N int64 // Total bytes written.
}

func (r *byteWriter) Write(p []byte) (n int, err error) {
	for len(p) > 0 {
		n := 16
		if n > len(p) {
			n = len(p)
		}
		for _, c := range p[:n] {
			//noerrcheck
			_, _ = fmt.Fprintf(r.w, "0x%02x,", c)
		}
		p = p[n:]
	}
	return n, nil
}
