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
	TYPE_SLICE
	TYPE_FUNCTION
	TYPE_VOID
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
	case TYPE_SLICE:
		return "slice"
	case TYPE_FUNCTION:
		return "function"
	case TYPE_VOID:
		return "void"
	}

	Unreachable("TypeEnum::String")
	return ""
}

type SliceType struct {
	of *Type
}

type FunctionType struct {
	parameters []*Type
	returnType *Type
}

func (t *FunctionType) Arity() int {
	return len(t.parameters)
}

type Type struct {
	kind  TypeEnum
	other interface{} // is either nil, or contains a `SliceType` or `FunctionType`
}

func (t Type) IsVoid() bool {
	return t.kind == TYPE_VOID
}

func (t Type) String() string {
	switch t.kind {
	case TYPE_I64, TYPE_U64, TYPE_BOOL, TYPE_STRING, TYPE_DOUBLE, TYPE_VOID:
		return t.kind.String()
	case TYPE_SLICE:
		other := t.other.(SliceType)
		if other.of.kind == TYPE_FUNCTION {
			/**
			 * if the `of` is a function, wrap the type in parentheses to resolve ambiguity
			 * e.g. `fn()int[]` is a function that returns an `int[]` while `(fn()int)[]` is a slice of functions
			 */
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

		fmt.Fprintf(&builder, ")%v", other.returnType)

		return builder.String()
	}

	Unreachable("Type::String")
	return ""
}

func TypesEqual(t1, t2 *Type) bool {
	if t1.kind != t2.kind {
		return false
	}

	switch t1.kind {
	case TYPE_I64, TYPE_U64, TYPE_BOOL, TYPE_STRING, TYPE_DOUBLE, TYPE_VOID:
		return true
	case TYPE_SLICE:
		other1 := t1.other.(SliceType)
		other2 := t2.other.(SliceType)
		return TypesEqual(other1.of, other2.of)
	case TYPE_FUNCTION:
		other1 := t1.other.(FunctionType)
		other2 := t2.other.(FunctionType)

		if other1.Arity() != other2.Arity() {
			return false
		}

		for i := range other1.parameters {
			if !TypesEqual(other1.parameters[i], other2.parameters[i]) {
				return false
			}
		}

		return TypesEqual(other1.returnType, other2.returnType)
	}

	Unreachable("types.go: TypesEqual()")
	return false
}

func SimpleType(typeEnum TypeEnum) *Type {
	return &Type{kind: typeEnum}
}

func SimpleFunction(returnType TypeEnum, parameters ...TypeEnum) FunctionType {
	parameterTypes := make([]*Type, len(parameters))
	for i := range parameters {
		parameterTypes[i] = SimpleType(parameters[i])
	}

	return FunctionType{
		returnType: SimpleType(returnType),
		parameters: parameterTypes,
	}
}
