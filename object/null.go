package object

// Null wraps nothing and implements our Object interface.
type Null struct{}

// Type returns the type of this object.
func (n *Null) Type() Type {
	return NULL_OBJ
}

// Inspect returns a string-representation of the given object.
func (n *Null) Inspect() string {
	return "null"
}

// GetMethod returns a method against the object.
// (Built-in methods only.)
func (n *Null) GetMethod(string) BuiltinFunction {
	return nil
}

// ToInterface converts this object to a go-interface, which will allow
// it to be used naturally in our sprintf/printf primitives.
func (n *Null) ToInterface() interface{} {
	return "<NULL>"
}

// JSON returns a json-friendly string
func (n *Null) JSON(indent bool) string {
	return n.Inspect()
}
