// The implementation of our string-object.

package object

import (
	"hash/fnv"
	"sort"
	"strings"
	"unicode/utf8"
)

// DocString wraps string and implements Object and Hashable interfaces.
type DocString struct {
	// Value holds the string value this object wraps.
	Value string

	// Offset holds our iteration-offset
	offset int
}

// Type returns the type of this object.
func (s *DocString) Type() Type {
	return DOCSTRING_OBJ
}

// Inspect returns an empty string;
// a docstring should only be printable when it's part of a block
func (s *DocString) Inspect() string {
	return ""
}

// HashKey returns a hash key for the given object.
func (s *DocString) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(s.Value))
	return HashKey{Type: s.Type(), Value: h.Sum64()}
}

// InvokeMethod invokes a method against the object.
// (Built-in methods only.)
func (s *DocString) InvokeMethod(method string, env Environment, args ...Object) Object {
	if method == "len" {
		return &Integer{Value: int64(utf8.RuneCountInString(s.Value))}
	}
	if method == "methods" {
		static := []string{"len", "methods"}
		dynamic := env.Names("docstring.")

		var names []string
		names = append(names, static...)

		for _, e := range dynamic {
			bits := strings.Split(e, ".")
			names = append(names, bits[1])
		}
		sort.Strings(names)

		result := make([]Object, len(names))
		for i, txt := range names {
			result[i] = &DocString{Value: txt}
		}
		return &Array{Elements: result}
	}
	return nil
}

// Reset implements the Iterable interface, and allows the contents
// of the string to be reset to allow re-iteration.
func (s *DocString) Reset() {
	s.offset = 0
}

// Next implements the Iterable interface, and allows the contents
// of our string to be iterated over.
func (s *DocString) Next() (Object, Object, bool) {

	if s.offset < utf8.RuneCountInString(s.Value) {
		s.offset++

		// Get the characters as an array of runes
		chars := []rune(s.Value)

		// Now index
		val := &DocString{Value: string(chars[s.offset-1])}

		return val, &Integer{Value: int64(s.offset - 1)}, true
	}

	return nil, &Integer{Value: 0}, false
}

// ToInterface converts this object to a go-interface, which will allow
// it to be used naturally in our sprintf/printf primitives.
//
// It might also be helpful for embedded users.
func (s *DocString) ToInterface() interface{} {
	return s.Value
}