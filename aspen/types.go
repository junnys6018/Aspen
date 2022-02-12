package main

import (
	"fmt"
	"strings"
)

type TypeEnum int

const (
	TYPE_I64 TypeEnum = iota
	TYPE_U64
	TYPE_BOOL
	TYPE_STRING
	TYPE_DOUBLE
	TYPE_NIL
	TYPE_ARRAY
	TYPE_SLICE
	TYPE_FUNCTION
)

func (t TypeEnum) IsNumeric() bool {
	return t == TYPE_I64 || t == TYPE_U64 || t == TYPE_DOUBLE
}

func (t TypeEnum) IsIntegral() bool {
	return t == TYPE_I64 || t == TYPE_U64
}

func (t TypeEnum) String() string {
	switch t {
	case TYPE_I64:
		return "i64"
	case TYPE_U64:
		return "u64"
	case TYPE_BOOL:
		return "bool"
	case TYPE_STRING:
		return "string"
	case TYPE_DOUBLE:
		return "double"
	case TYPE_NIL:
		return "nil"
	case TYPE_ARRAY:
		return "array"
	case TYPE_SLICE:
		return "slice"
	case TYPE_FUNCTION:
		return "function"
	default:
		panic("TypeEnum::String unknown type")
	}
}

type ArrayType struct {
	count int
	of    Type
}

type SliceType struct {
	of Type
}

type FunctionType struct {
	parameters []Type
	returnType Type
}

func (t *FunctionType) Arity() int {
	return len(t.parameters)
}

type Type struct {
	kind  TypeEnum
	other interface{} // is either nil, or contains an `ArrayType`, `SliceType` or `FunctionType`
}

func (t Type) String() string {
	switch t.kind {
	case TYPE_I64, TYPE_U64, TYPE_BOOL, TYPE_STRING, TYPE_DOUBLE, TYPE_NIL:
		return t.kind.String()
	case TYPE_ARRAY:
		other := t.other.(ArrayType)
		if other.of.kind == TYPE_FUNCTION {
			/**
			 * if the `of` is a function, wrap the type in parentheses to resolve ambiguity
			 * e.g. `fn()int[]` is a function that returns an `int[]` while `(fn()int)[]` is a slice of functions
			 */
			return fmt.Sprintf("(%v)[%d]", other.of, other.count)
		} else {
			return fmt.Sprintf("%v[%d]", other.of, other.count)
		}
	case TYPE_SLICE:
		other := t.other.(SliceType)
		if other.of.kind == TYPE_FUNCTION {
			return fmt.Sprintf("(%v)[]", other.of)
		} else {
			return fmt.Sprintf("%v[]", other.of)
		}
	case TYPE_FUNCTION:
		other := t.other.(FunctionType)
		builder := strings.Builder{}

		fmt.Fprintf(&builder, "fn(")
		for i, parameter := range other.parameters {
			fmt.Fprintf(&builder, "%v", parameter)
			if i != len(other.parameters)-1 {
				builder.WriteString(", ")
			}
		}
		if other.returnType.kind == TYPE_NIL {
			fmt.Fprintf(&builder, ")") // function returns nothing, omit return type
		} else {
			fmt.Fprintf(&builder, ")%v", other.returnType)
		}

		return builder.String()
	default:
		panic("TypeEnum::String unknown type")
	}
}

func SimpleType(typeEnum TypeEnum) Type {
	return Type{kind: typeEnum}
}
