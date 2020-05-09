package strings

import (
	"math/rand"
	"unicode"
)

// Type represents bit mask of rune types according to Unicode.
type Type int

// Is reports whether the type includes a given one.
func (t Type) Is(t2 Type) bool { return t&t2 == t2 }

// Equal reports whether the type is equal to a given one.
func (t Type) Equal(t2 Type) bool { return t == t2 }

const (
	Letter Type = 1 << iota
	Number
	Symbol
	Graphic
	Digit
	Mark
	Control
	Lower
	Upper
	Print
	Punct
	Space
	Title
)

type is func(typ Type, r rune) Type

func foreach(is ...is) is {
	return func(typ Type, r rune) Type {
		for _, i := range is {
			typ = i(typ, r)
		}
		return typ
	}
}

func when(is func(rune) bool, then Type) is {
	return func(current Type, r rune) Type {
		if is(r) {
			current |= then
		}
		return current
	}
}

var rtype = foreach(
	when(unicode.IsLetter, Letter),
	when(unicode.IsNumber, Number),
	when(unicode.IsSymbol, Symbol),
	when(unicode.IsGraphic, Graphic),
	when(unicode.IsDigit, Digit),
	when(unicode.IsMark, Mark),
	when(unicode.IsControl, Control),
	when(unicode.IsLower, Lower),
	when(unicode.IsUpper, Upper),
	when(unicode.IsPrint, Print),
	when(unicode.IsPunct, Punct),
	when(unicode.IsSpace, Space),
	when(unicode.IsTitle, Title),
)

// Rune represents std rune wrapper with it Unicode type.
type Rune struct {
	rune
	typ Type
}

// Returns Rune representation of std rune.
func NewRune(r rune) Rune {
	return Rune{
		rune: r,
		typ:  rtype(0, r),
	}
}

// Equal reports whether the rune is equal to given one.
func (r Rune) Equal(r2 Rune) bool { return r.rune == r2.rune }

// Same reports whether the rune is equal to given one regardless to its case.
func (r Rune) Same(r2 Rune) bool { return unicode.ToUpper(r.rune) == unicode.ToUpper(r2.rune) }

// Rune returns std rune.
func (r Rune) Rune() rune { return r.rune }

// OfType reports whether the rune's type is equal to given one.
func (r Rune) OfType(r2 Rune) bool { return r.typ.Equal(r2.typ) }

// Is reports whether the rune includes type of given one.
func (r Rune) Is(typ Type) bool { return r.typ.Is(typ) }

// Type returns rune's type.
func (r Rune) Type() Type { return r.typ }

// String returns string representation of rune.
func (r Rune) String() string { return string(r.rune) }

// String represents alias on array of Runes.
type String []Rune

// New returns String representation of std string.
func New(s string) String {
	runes := []rune(s)
	str := make(String, len(runes))
	for i, r := range runes {
		str[i] = NewRune(r)
	}
	return str
}

// Len returns string length.
func (s String) Len() int { return len(s) }

// Empty reports whether the String is of 0 length.
func (s String) Empty() bool { return len(s) == 0 }

// Bytes returns byte slice representation of String.
func (s String) Bytes() []byte { return []byte(s.String()) }

// Shuffle pseudo-randomizes the order of runes
// and returns a copy of original String.
func (s String) Shuffle(rnd *rand.Rand) String {
	rnd.Shuffle(len(s), func(i, j int) { s[i], s[j] = s[j], s[i] })
	return s.Copy()
}

// Copy returns copy of original String.
func (s String) Copy() String {
	copied := make(String, len(s))
	copy(copied, s)
	return copied
}

// Remove removes i element from String if i
// is in range of [0:len(s)-1] and returns a copy
// of original String.
func (s String) Remove(i int) (String, bool) {
	if i < 0 || i >= len(s) {
		return s, false
	}
	return append(s[:i], s[i+1:]...), true
}

// Pick returns i Rune of String if i is
// in range of [0:len(s)-1].
func (s String) Pick(i int) (r Rune, picked bool) {
	if i < 0 || i >= len(s) {
		return r, false
	}
	return s[i], true
}

// Last returns last Rune of String.
func (s String) Last() (Rune, bool) { return s.Pick(len(s) - 1) }

// String returns string representation of String.
func (s String) String() string {
	result := make([]rune, len(s))
	for i, r := range s {
		result[i] = r.rune
	}
	return string(result)
}
