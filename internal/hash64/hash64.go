package hash64

import (
	"crypto/sha512"
	"hash"
	"hash/fnv"
	"io"
)

type (
	sum interface {
		// Sum appends the current hash to b and returns the resulting slice.
		// It does not change the underlying hash state.
		Sum(b []byte) []byte

		// Size returns the number of bytes Sum will return.
		Size() int

		// BlockSize returns the hash's underlying block size.
		// The Write method must be able to accept any amount
		// of data, but it may operate more efficiently if all writes
		// are a multiple of the block size.
		BlockSize() int
	}
	sum64 interface {
		Sum64() uint64
	}
)

type hash64 struct {
	io.Writer
	sum
	sum64
	reset []interface{ Reset() }
}

// New returns hash.Hash64 implementation
// where Sum implements by sha-512
// and Sum64 implements by fnv64.
func New() hash.Hash64 {
	sum64 := fnv.New64()
	sum := sha512.New()
	return &hash64{
		Writer: io.MultiWriter(sum64, sum),
		sum:    sum,
		sum64:  sum64,
		reset:  []interface{ Reset() }{sum64, sum},
	}
}

func (h *hash64) Reset() {
	for _, r := range h.reset {
		r.Reset()
	}
}
