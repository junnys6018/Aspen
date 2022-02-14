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

	array := Type{kind: TYPE_SLICE, other: SliceType{of: SimpleType(TYPE_I64)}}
	expect(array.String(), "i64[]")

	function := Type{kind: TYPE_FUNCTION, other: FunctionType{parameters: []*Type{SimpleType(TYPE_I64), SimpleType(TYPE_STRING)}, returnType: SimpleType(TYPE_BOOL)}}
	expect(function.String(), "fn(i64, string)bool")

	ambiguous1 := Type{kind: TYPE_SLICE, other: SliceType{of: &Type{kind: TYPE_FUNCTION, other: FunctionType{parameters: []*Type{}, returnType: SimpleType(TYPE_I64)}}}}
	expect(ambiguous1.String(), "(fn()i64)[]")

	ambiguous2 := Type{kind: TYPE_FUNCTION, other: FunctionType{parameters: []*Type{}, returnType: &Type{kind: TYPE_SLICE, other: SliceType{of: SimpleType(TYPE_I64)}}}}
	expect(ambiguous2.String(), "fn()i64[]")

	noReturn := Type{kind: TYPE_FUNCTION, other: FunctionType{parameters: []*Type{}, returnType: SimpleType(TYPE_VOID)}}
	expect(noReturn.String(), "fn()void")
}

func TestTypesEqual(t *testing.T) {
	var t1, t2 *Type

	expectEqual := func() {
		if !TypesEqual(t1, t2) {
			t.Errorf("expected %v to equal %v", t1, t2)
		}
	}

	expectNotEqual := func() {
		if TypesEqual(t1, t2) {
			t.Errorf("expected %v to not equal %v", t1, t2)
		}
	}

	t1 = SimpleType(TYPE_BOOL)
	t2 = SimpleType(TYPE_DOUBLE)
	expectNotEqual()

	t1 = SimpleType(TYPE_I64)
	t2 = SimpleType(TYPE_I64)
	expectEqual()

	t1 = &Type{kind: TYPE_SLICE, other: SliceType{of: SimpleType(TYPE_I64)}}
	t2 = &Type{kind: TYPE_SLICE, other: SliceType{of: SimpleType(TYPE_I64)}}
	expectEqual()

	t1 = &Type{kind: TYPE_SLICE, other: SliceType{of: SimpleType(TYPE_I64)}}
	t2 = &Type{kind: TYPE_SLICE, other: SliceType{of: SimpleType(TYPE_BOOL)}}
	expectNotEqual()

	t1 = &Type{kind: TYPE_SLICE, other: SliceType{of: &Type{kind: TYPE_SLICE, other: SliceType{of: SimpleType(TYPE_I64)}}}}
	t2 = &Type{kind: TYPE_SLICE, other: SliceType{of: SimpleType(TYPE_I64)}}
	expectNotEqual()

	t1 = &Type{kind: TYPE_SLICE, other: SliceType{of: &Type{kind: TYPE_SLICE, other: SliceType{of: SimpleType(TYPE_I64)}}}}
	t2 = &Type{kind: TYPE_SLICE, other: SliceType{of: &Type{kind: TYPE_SLICE, other: SliceType{of: SimpleType(TYPE_I64)}}}}
	expectEqual()

	t1 = &Type{kind: TYPE_FUNCTION, other: FunctionType{parameters: []*Type{}, returnType: SimpleType(TYPE_BOOL)}}
	t2 = &Type{kind: TYPE_FUNCTION, other: FunctionType{parameters: []*Type{}, returnType: SimpleType(TYPE_BOOL)}}
	expectEqual()

	t1 = &Type{kind: TYPE_FUNCTION, other: FunctionType{parameters: []*Type{SimpleType(TYPE_BOOL)}, returnType: SimpleType(TYPE_BOOL)}}
	t2 = &Type{kind: TYPE_FUNCTION, other: FunctionType{parameters: []*Type{SimpleType(TYPE_BOOL), SimpleType(TYPE_BOOL)}, returnType: SimpleType(TYPE_BOOL)}}
	expectNotEqual()

	t1 = &Type{kind: TYPE_FUNCTION, other: FunctionType{parameters: []*Type{SimpleType(TYPE_BOOL)}, returnType: SimpleType(TYPE_I64)}}
	t2 = &Type{kind: TYPE_FUNCTION, other: FunctionType{parameters: []*Type{SimpleType(TYPE_BOOL)}, returnType: SimpleType(TYPE_BOOL)}}
	expectNotEqual()

	t1 = &Type{kind: TYPE_FUNCTION, other: FunctionType{parameters: []*Type{SimpleType(TYPE_I64)}, returnType: SimpleType(TYPE_I64)}}
	t2 = &Type{kind: TYPE_FUNCTION, other: FunctionType{parameters: []*Type{SimpleType(TYPE_BOOL)}, returnType: SimpleType(TYPE_I64)}}
	expectNotEqual()
}
