package main

import (
	"testing"
)

func TestTypesStringFunction(t *testing.T) {
	expect := func(s1, s2 string) {
		if s1 != s2 {
			t.Errorf("expected \"%s\", got \"%s\"", s2, s1)
		}
	}

	simple := SimpleType(TYPE_DOUBLE)
	expect(simple.String(), "double")

	array := Type{kind: TYPE_ARRAY, other: ArrayType{of: SimpleType(TYPE_I64), count: 42}}
	expect(array.String(), "i64[42]")

	function := Type{kind: TYPE_FUNCTION, other: FunctionType{parameters: []Type{SimpleType(TYPE_I64), SimpleType(TYPE_STRING)}, returnType: SimpleType(TYPE_BOOL)}}
	expect(function.String(), "fn(i64, string)bool")

	ambiguous1 := Type{kind: TYPE_SLICE, other: SliceType{of: Type{kind: TYPE_FUNCTION, other: FunctionType{parameters: []Type{}, returnType: SimpleType(TYPE_I64)}}}}
	expect(ambiguous1.String(), "(fn()i64)[]")

	ambiguous2 := Type{kind: TYPE_FUNCTION, other: FunctionType{parameters: []Type{}, returnType: Type{kind: TYPE_SLICE, other: SliceType{of: SimpleType(TYPE_I64)}}}}
	expect(ambiguous2.String(), "fn()i64[]")

	noReturn := Type{kind: TYPE_FUNCTION, other: FunctionType{parameters: []Type{}, returnType: SimpleType(TYPE_NIL)}}
	expect(noReturn.String(), "fn()")
}
