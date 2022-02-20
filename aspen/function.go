package main

import "fmt"

type AspenFunction interface {
	Arity() int
	Call(interpreter *Interpreter, args []interface{}) interface{}
	String() string
}

type NativeFunction struct {
	atype FunctionType
	impl  func([]interface{}) interface{}
}

func (f *NativeFunction) Arity() int {
	return f.atype.Arity()
}

func (f *NativeFunction) Call(interpreter *Interpreter, args []interface{}) interface{} {
	return f.impl(args)
}

func (f *NativeFunction) String() string {
	return "<native fn>"
}

type UserFunction struct {
	declaration *FunctionStatement
	closure     Environment
}

type ReturnValue struct {
	value interface{}
}

func (f *UserFunction) Arity() int {
	return f.declaration.atype.Arity()
}

func (f *UserFunction) Call(interpreter *Interpreter, args []interface{}) (ret interface{}) {
	defer func() {
		if r := recover(); r != nil {
			switch v := r.(type) {
			case ReturnValue:
				// recover from any calls to panic with an argument of type `ReturnValue`
				// override the return value
				ret = v.value
			default:
				// else re-panic
				panic(v)
			}
		}
	}()

	environment := NewEnvironment(&f.closure)

	for i := range f.declaration.parameters {
		environment.Define(f.declaration.parameters[i].String(), args[i])
	}

	interpreter.ExecuteBlock(f.declaration.body, environment)

	return nil
}

func (f *UserFunction) String() string {
	return fmt.Sprintf("<fn %v>", f.declaration.name)
}

var NativeFunctions = make(map[string]*NativeFunction)

func DefineNativeFunction(atype FunctionType, name string, impl func([]interface{}) interface{}) {
	NativeFunctions[name] = &NativeFunction{atype: atype, impl: impl}
}
