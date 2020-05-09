package swca1_test

import (
	"fmt"

	"github.com/unqnown/swca1"
)

func Example() {
	h := swca1.New(
		swca1.Size(swca1.Max),
		swca1.Alphabet(swca1.NULS),
		swca1.Unique,
		swca1.NoTypeRepetition,
		swca1.NoCharacterRepetition,
	)

	_, _ = h.Write([]byte("salt"))

	_, _ = h.Write([]byte("hint"))

	fmt.Printf("%s", h.Sum(nil))
	// output: e@c*x6d-q#J8U&F^%H+t7Y!L1P4M5I$g0S2A?R3o9n
}
