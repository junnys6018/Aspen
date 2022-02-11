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
	count    int
	category TypeEnum
}

type SliceType struct {
	category TypeEnum
}

type FunctionType struct {
	parameters []TypeEnum
	returnType TypeEnum
}

func (t *FunctionType) Arity() int {
	return len(t.parameters)
}

type Type struct {
	category TypeEnum
	other    interface{}
}

func (t Type) String() string {
	switch t.category {
	case TYPE_I64, TYPE_U64, TYPE_BOOL, TYPE_STRING, TYPE_DOUBLE, TYPE_NIL:
		return t.category.String()
	case TYPE_ARRAY:
		other := t.other.(ArrayType)
		return fmt.Sprintf("%v[%d]", other.category, other.count)
	case TYPE_SLICE:
		other := t.other.(SliceType)
		return fmt.Sprintf("%v[]", other.category)
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
		fmt.Fprintf(&builder, ")%v", other.returnType)

		return builder.String()
	default:
		panic("TypeEnum::String unknown type")
	}
}

func SimpleType(typeEnum TypeEnum) Type {
	return Type{category: typeEnum}
}
