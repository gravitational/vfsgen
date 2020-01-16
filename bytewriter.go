package vfsgen

import (
	"fmt"
	"io"
	"strings"
)

// byteWriter encodes the input as a []byte literal.
// It outputs the byte slice as gofmt would format it without requiring gofmt.
type byteWriter struct {
	w     io.Writer
	ntabs int
	err   error
}

func (r *byteWriter) Write(p []byte) (n int, err error) {
	const numBytesInRow = 16
	for len(p) > 0 {
		numBytes := numBytesInRow
		if numBytes > len(p) {
			// Write the rest of bytes
			numBytes = len(p)
		}
		r.indent()
		r.write(p[0])
		// Have any bytes to output?
		if numBytes > 1 {
			for _, c := range p[1:numBytes] {
				r.writeAdditional(c)
			}
		}
		if numBytes == numBytesInRow {
			// Start a new line
			r.writeEOL()
		}
		n += numBytes
		p = p[numBytes:]
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

func (r *byteWriter) writeEOL() error {
	if r.err != nil {
		return r.err
	}
	_, r.err = fmt.Fprint(r.w, "\n")
	return r.err
}
