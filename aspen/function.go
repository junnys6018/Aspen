package main

type AspenFunction interface {
	Arity() int
	Call(arguments []interface{}) interface{}
	String() string
}

type NativeFunction struct {
	atype FunctionType
	impl  func([]interface{}) interface{}
}

func (f *NativeFunction) Arity() int {
	return f.atype.Arity()
}

func (f *NativeFunction) Call(args []interface{}) interface{} {
	return f.impl(args)
}

func (f *NativeFunction) String() string {
	return "<native fn>"
}

var NativeFunctions = make(map[string]*NativeFunction)

func DefineNativeFunction(atype FunctionType, name string, impl func([]interface{}) interface{}) {
	NativeFunctions[name] = &NativeFunction{atype: atype, impl: impl}
}
