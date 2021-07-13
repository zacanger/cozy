// The implementation of our regular-expression object.

package object

// Regexp wraps regular-expressions and implements the Object interface.
type Regexp struct {
	// Value holds the string value this object wraps.
	Value string

	// Flags holds the flags for the object
	Flags string
}

// Type returns the type of this object.
func (r *Regexp) Type() Type {
	return REGEXP_OBJ
}

// Inspect returns a string-representation of the given object.
func (r *Regexp) Inspect() string {
	return r.Value
}

// GetMethod returns a method against the object.
// (Built-in methods only.)
func (r *Regexp) GetMethod(string) BuiltinFunction {
	return nil
}

// ToInterface converts this object to a go-interface, which will allow
// it to be used naturally in our sprintf/printf primitives.
//
// It might also be helpful for embedded users.
func (r *Regexp) ToInterface() interface{} {
	return "<REGEXP>"
}
