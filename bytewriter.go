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
		s := 16
		if s > len(p) {
			s = len(p)
		}
		for _, c := range p[:s] {
			//noerrcheck
			_, err = fmt.Fprintf(r.w, "0x%02x,", c)
			if err != nil {
				return n, err
			}
		}
		n += s
		p = p[s:]
	}
	return n, nil
}
