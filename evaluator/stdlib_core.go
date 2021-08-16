package evaluator

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/zacanger/cozy/object"
	"github.com/zacanger/cozy/utils"
)

// convert a string to a float
func floatFn(args ...object.Object) object.Object {
	if len(args) != 1 {
		return NewError("wrong number of arguments. got=%d, want=1",
			len(args))
	}
	switch args[0].(type) {
	case *object.String:
		input := args[0].(*object.String).Value
		i, err := strconv.Atoi(input)
		if err == nil {
			return &object.Float{Value: float64(i)}
		}
		return NewError("Converting string '%s' to float failed %s", input, err.Error())

	case *object.Boolean:
		input := args[0].(*object.Boolean).Value
		if input {
			return &object.Float{Value: float64(1)}

		}
		return &object.Float{Value: float64(0)}
	case *object.Float:
		// noop
		return args[0]
	case *object.Integer:
		input := args[0].(*object.Integer).Value
		return &object.Float{Value: float64(input)}
	default:
		return NewError("argument to `float` not supported, got=%s",
			args[0].Type())
	}
}

// convert a double/string to an int
func intFn(args ...object.Object) object.Object {
	if len(args) != 1 {
		return NewError("wrong number of arguments. got=%d, want=1",
			len(args))
	}
	switch args[0].(type) {
	case *object.String:
		input := args[0].(*object.String).Value
		i, err := strconv.Atoi(input)
		if err == nil {
			return &object.Integer{Value: int64(i)}
		}
		return NewError("Converting string '%s' to int failed %s", input, err.Error())

	case *object.Boolean:
		input := args[0].(*object.Boolean).Value
		if input {
			return &object.Integer{Value: 1}

		}
		return &object.Integer{Value: 0}
	case *object.Integer:
		// noop
		return args[0]
	case *object.Float:
		input := args[0].(*object.Float).Value
		return &object.Integer{Value: int64(input)}
	default:
		return NewError("argument to `int` not supported, got=%s",
			args[0].Type())
	}
}

// length of item
func lenFn(args ...object.Object) object.Object {
	if len(args) != 1 {
		return NewError("wrong number of arguments. got=%d, want=1",
			len(args))
	}
	switch arg := args[0].(type) {
	case *object.String:
		return &object.Integer{Value: int64(utf8.RuneCountInString(arg.Value))}
	case *object.DocString:
		return &object.Integer{Value: int64(utf8.RuneCountInString(arg.Value))}
	case *object.Array:
		return &object.Integer{Value: int64(len(arg.Elements))}
	case *object.Null:
		return &object.Integer{Value: 0}
	case *object.Hash:
		return &object.Integer{Value: int64(len(arg.Pairs))}
	default:
		return NewError("argument to `len` not supported, got=%s",
			args[0].Type())
	}
}

// regular expression match
func matchFn(args ...object.Object) object.Object {
	if len(args) != 2 {
		return NewError("wrong number of arguments. got=%d, want=2",
			len(args))
	}

	if args[0].Type() != object.STRING_OBJ {
		return NewError("argument to `match` must be STRING, got %s",
			args[0].Type())
	}
	if args[1].Type() != object.STRING_OBJ {
		return NewError("argument to `match` must be STRING, got %s",
			args[1].Type())
	}

	// Compile and match
	reg := regexp.MustCompile(args[0].(*object.String).Value)
	res := reg.FindStringSubmatch(args[1].(*object.String).Value)

	if len(res) > 0 {
		newArray := make([]object.Object, len(res))

		// If we get a match then the output is an array
		// First entry is the match, any additional parts
		// are the capture-groups.
		if len(res) > 1 {
			for i, v := range res {
				newArray[i] = &object.String{Value: v}
			}
		}

		return &object.Array{Elements: newArray}
	}

	// No match
	return &object.Array{Elements: make([]object.Object, 0)}
}

// output a string to stdout
func printFn(args ...object.Object) object.Object {
	for _, arg := range args {
		fmt.Print(arg.Inspect() + " ")
	}
	fmt.Print("\n")
	return NULL
}

func strFn(args ...object.Object) object.Object {
	if len(args) != 1 {
		return NewError("wrong number of arguments. got=%d, want=1",
			len(args))
	}

	out := args[0].Inspect()
	return &object.String{Value: out}
}

// type of an item
func typeFn(args ...object.Object) object.Object {
	if len(args) != 1 {
		return NewError("wrong number of arguments. got=%d, want=1",
			len(args))
	}
	return &object.String{Value: strings.ToLower(string(args[0].Type()))}
}

// error
func errorFn(args ...object.Object) object.Object {
	if len(args) != 1 {
		return NewError("wrong number of arguments. got=%d, want=1",
			len(args))
	}
	switch t := args[0].(type) {
	case *object.String:
		return &object.Error{Message: t.Value, BuiltinCall: true}
	case *object.Hash:
		msgStr := &object.String{Value: "message"}
		codeStr := &object.String{Value: "code"}
		dataStr := &object.String{Value: "data"}
		msg := t.Pairs[msgStr.HashKey()].Value
		code := t.Pairs[codeStr.HashKey()].Value
		data := t.Pairs[dataStr.HashKey()].Value
		e := &object.Error{BuiltinCall: true}
		if msg != nil {
			switch m := msg.(type) {
			case *object.String:
				e.Message = m.Value
			default:
				return NewError("error.message should be string!")
			}
		}
		if code != nil {
			switch c := code.(type) {
			case *object.Integer:
				cc := int(c.Value)
				e.Code = &cc
			default:
				return NewError("error.code should be integer!")
			}
		}
		if data != nil {
			e.Data = data.Json()
		}
		return e
	default:
		return NewError("error() expected a string or hash!")
	}
}

// panic
func panicFn(args ...object.Object) object.Object {
	switch e := args[0].(type) {
	case *object.Error:
		c := 1
		fmt.Println(e.Message)
		if e.Code != nil {
			c = int(*e.Code)
		}
		utils.ExitConditionally(c)
	default:
		return NewError("panic expected an error!")
	}
	return NULL
}

func init() {
	RegisterBuiltin("int",
		func(env *object.Environment, args ...object.Object) object.Object {
			return intFn(args...)
		})
	RegisterBuiltin("float",
		func(env *object.Environment, args ...object.Object) object.Object {
			return floatFn(args...)
		})
	RegisterBuiltin("len",
		func(env *object.Environment, args ...object.Object) object.Object {
			return lenFn(args...)
		})
	RegisterBuiltin("match",
		func(env *object.Environment, args ...object.Object) object.Object {
			return matchFn(args...)
		})
	RegisterBuiltin("print",
		func(env *object.Environment, args ...object.Object) object.Object {
			return printFn(args...)
		})
	RegisterBuiltin("string",
		func(env *object.Environment, args ...object.Object) object.Object {
			return strFn(args...)
		})
	RegisterBuiltin("type",
		func(env *object.Environment, args ...object.Object) object.Object {
			return typeFn(args...)
		})
	RegisterBuiltin("error",
		func(env *object.Environment, args ...object.Object) object.Object {
			return errorFn(args...)
		})
	RegisterBuiltin("panic",
		func(env *object.Environment, args ...object.Object) object.Object {
			return panicFn(args...)
		})
}
