package randomer

import (
	"crypto/rand"
	"io"
)

type randomer struct{}

// New returns io.Reader which reads randome data.
func New() io.Reader {
	return &randomer{}
}

var _ io.Reader = (*randomer)(nil)

func (r *randomer) Read(p []byte) (int, error) {
	return rand.Read(p)
}
