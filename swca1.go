package swca1

import (
	"errors"
	"fmt"
	"hash"
	"io"
	"math"
	"math/rand"

	"github.com/unqnown/swca1/internal/hash64"
	"github.com/unqnown/swca1/internal/strings"

	stdstrings "strings"
)

var vocabulary = map[rune]string{
	'1': "0123456789",
	'n': "0123456789",
	'A': "ABCDEFGHIJKLMNOPQRSTUVWXYZ",
	'u': "ABCDEFGHIJKLMNOPQRSTUVWXYZ",
	'a': "abcdefghijklmnopqrstuvwxyz",
	'l': "abcdefghijklmnopqrstuvwxyz",
	'@': "?!@#$%^&*-+",
	's': "?!@#$%^&*-+",
}

// Options represents hash generation configuration.
type Options struct {
	size     int
	alphabet string
	complex  complexity
}

type complexity func(strings.String, strings.Rune) error

func wrap(v complexity, w ...func(complexity) complexity) complexity {
	if len(w) == 0 {
		return v
	}
	return w[0](wrap(v, w[1:cap(w)]...))
}

// Alphabet allows to customize hash alphabet in "wkt|alphabet" format
// where "wkt" is a well known tokens:
// 	- n(1) - numbers;
//  - u(A) - uppercase letters;
//	- l(a) - lowercase letters;
//	- s(@) - symbols.
// Rest of tokens string after "|" is a custom runes which will
// append to result alphabet.
func Alphabet(tokens string) Option {
	return func(o *Options) {
		if len(tokens) == 0 {
			return
		}
		o.alphabet = ""
		splits := stdstrings.SplitN(tokens, "|", 2)
		for _, t := range splits[0] {
			s, found := vocabulary[t]
			if !found {
				continue
			}
			o.alphabet += s
		}
		if len(splits) == 2 {
			o.alphabet += splits[1]
		}
	}
}

// Size allows to specify hash size. Allowed hash size is in range [0:64]
// where 0 is reserved for specifying max hash size within
// required complexity.
func Size(size int) Option {
	return func(o *Options) {
		if size >= 0 {
			o.size = size
		}
	}
}

// Unique ensures uniqueness of each character in hash.
func Unique(o *Options) { o.complex = wrap(o.complex, unique) }

// NoTypeRepetition guarantees any pair of adjacent characters
// in the hash will have a different type.
func NoTypeRepetition(o *Options) { o.complex = wrap(o.complex, noTypeRepetition) }

// NoCharacterRepetition guarantees that in any pair of adjacent
// characters there will be no more than one letter, regardless of case.
func NoCharacterRepetition(o *Options) { o.complex = wrap(o.complex, noCharacterRepetition) }

func unique(v complexity) complexity {
	return func(h strings.String, r strings.Rune) error {
		for _, c := range h {
			if c.Same(r) {
				return fmt.Errorf("swca1: incoming rune %q is already in hash", r)
			}
		}
		return v(h, r)
	}
}

func noTypeRepetition(v complexity) complexity {
	return func(h strings.String, r strings.Rune) error {
		prev, has := h.Last()
		if !has {
			return nil
		}
		if prev.OfType(r) {
			return fmt.Errorf("swca1: repeated type in %q %q pair", prev, r)
		}
		return v(h, r)
	}
}

func noCharacterRepetition(v complexity) complexity {
	return func(h strings.String, r strings.Rune) error {
		prev, has := h.Last()
		if !has {
			return nil
		}
		if prev.Is(strings.Letter) && r.Is(strings.Letter) {
			return fmt.Errorf("swca1: repeated character in %q %q pair", prev, r)
		}
		return v(h, r)
	}
}

// Option is an type of go idiomatic options representation.
type Option func(*Options)

type digest struct {
	opts Options
	hash.Hash64
	w   io.Writer
	rnd *rand.Rand

	hash strings.String
}

// ErrUnreachableComplexity means that there is no way to generate a hash
// of such size with requested complexity.
var ErrUnreachableComplexity = errors.New("swca1: unreachable complexity")

const (
	// NULS represents default alphabet's tokens.
	NULS = "nuls"
	// ABC represents default alphabet.
	ABC = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz?!@#$%^&*-+"
	// Enough is an optimal hash size for default alphabet.
	Enough = 20
	// Max is reserved identifier for maximum hash size within required complexity.
	Max = 0
)

// New returns sha-512 based hash.Hash with custom alphabet.
func New(opts ...Option) hash.Hash {
	d := &digest{
		opts: Options{
			alphabet: ABC,
			complex: func(strings.String, strings.Rune) error {
				return nil
			},
			size: Enough,
		},
		rnd:    rand.New(rand.NewSource(1)),
		Hash64: hash64.New(),
	}
	for _, opt := range opts {
		opt(&d.opts)
	}
	return d
}

func (d *digest) Write(data []byte) (n int, err error) {
	if n, err = d.Hash64.Write(data); err != nil {
		return n, err
	}

	sum := d.Hash64.Sum64()
	picks := btoi(d.Hash64.Sum(nil))

	switch {
	case d.opts.size == Max:
		//
	case d.opts.size < d.Hash64.Size():
		picks = picks[:d.opts.size]
	}

	abc := strings.New(d.opts.alphabet)
	d.hash = make(strings.String, 0, len(picks))

	var next strings.Rune

	var entropy int
	for _, p := range picks {
		d.rnd.Seed(int64(sum) + int64(p) + int64(entropy))
		src := abc.Shuffle(d.rnd)
		for attempt := 0; ; attempt++ {
			p = sin(p+attempt, src.Len()-1)
			entropy += p
			src, next = pick(src, p)
			if src.Empty() {
				if d.opts.size == Max {
					return len(data), nil
				}
				return 0, ErrUnreachableComplexity
			}
			if err := d.opts.complex(d.hash, next); err == nil {
				break
			}
		}
		d.hash = append(d.hash, next)
	}

	return len(data), nil
}

func (d *digest) Sum(b []byte) []byte { return append(b, d.hash.Bytes()...) }

func (d *digest) Reset() {
	d.hash = nil
	d.rnd.Seed(1)
	d.Hash64.Reset()
}
func (d *digest) Size() int {
	if d.opts.size > 0 {
		return d.opts.size
	}
	return d.Hash64.Size()
}
func (d *digest) BlockSize() int { return 1 }

func pick(abc strings.String, i int) (strings.String, strings.Rune) {
	pick, _ := abc.Pick(i)
	abc, _ = abc.Remove(i)
	return abc, pick
}

// sin calculates sin(x) within [0:max] range.
func sin(x, max int) int {
	y := math.Sin(float64(x))
	return int(y*float64(max)/2 + float64(max)/2)
}

func btoi(bytes []byte) []int {
	split := make([]int, len(bytes))
	for i, b := range bytes {
		split[i] = int(b)
	}
	return split
}
