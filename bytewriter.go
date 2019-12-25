package vfsgen

import (
	"fmt"
	"io"
	"strings"
)

// byteWriter encodes the input as a []byte literal.
type byteWriter struct {
	w     io.Writer
	ntabs int
	err   error
}

func (r *byteWriter) Write(p []byte) (n int, err error) {
	for len(p) > 0 {
		s := 16
		if s > len(p) {
			s = len(p)
		}
		r.indent()
		r.write(p[0])
		if s > 1 {
			for _, c := range p[1:s] {
				r.writeAdditional(c)
			}
		}
		if len(p) > 0 && s == 16 {
			r.writeEol()
		}
		n += s
		p = p[s:]
	}
	return n, r.err
}

func (r *byteWriter) indent() error {
	if r.err != nil {
		return r.err
	}
	_, r.err = fmt.Fprint(r.w, strings.Repeat("\t", r.ntabs))
	return r.err
}

func (r *byteWriter) write(c byte) error {
	if r.err != nil {
		return r.err
	}
	_, r.err = fmt.Fprintf(r.w, "0x%02x,", c)
	return r.err
}

func (r *byteWriter) writeAdditional(c byte) error {
	if r.err != nil {
		return r.err
	}
	_, r.err = fmt.Fprintf(r.w, " 0x%02x,", c)
	return r.err
}

func (r *byteWriter) writeEol() error {
	if r.err != nil {
		return r.err
	}
	_, r.err = fmt.Fprint(r.w, "\n")
	return r.err
}
