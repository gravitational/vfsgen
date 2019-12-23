package vfsgen

import (
	"fmt"
	"io"
)

// byteWriter encodes the input as a []byte literal.
// It tracks the total number of bytes written.
type byteWriter struct {
	w   io.Writer
	err error
	N   int64 // Total bytes written.
}

func (r *byteWriter) Write(p []byte) (n int, err error) {
	for len(p) > 0 {
		s := 16
		if s > len(p) {
			s = len(p)
		}
		for _, c := range p[:s] {
			r.write(c)
		}
		r.writeEol()
		n += s
		p = p[s:]
	}
	return n, r.err
}

func (r *byteWriter) write(c byte) error {
	if r.err != nil {
		return r.err
	}
	_, r.err = fmt.Fprintf(r.w, "0x%02x,", c)
	return r.err
}

func (r *byteWriter) writeEol() error {
	if r.err != nil {
		return r.err
	}
	_, r.err = fmt.Fprint(r.w, "\n")
	return r.err
}
