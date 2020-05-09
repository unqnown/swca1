### swca1

swca1 is zero dependency sha-512 based hash implementation with a custom alphabet.

### features

- hash alphabet customization in "wkt|alphabet" format, where "wkt" is a well-known tokens:
    * n or 1 - numbers;
    * u or A - uppercase letters;
    * l or a - lowercase letters;
    * s or @ - symbols;
    * rest of tokens string after "|" is a custom runes which will append to result alphabet.
- hash size specification in range [0:64] where 0 is reserved for specifying max hash size within required complexity;
- hash complexity specification:
    * unique - ensures uniqueness of each character in hash;
    * no type repetition - guarantees any pair of adjacent characters in the hash will have a different type;
    * no characters repetition - guarantees that in any pair of adjacent characters there will be no more than one letter, regardless of case. 

### installation

Standard `go get`:

```shell script
go get github.com/unqnown/swca1
````

### usage & example

```go
package test

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
```
