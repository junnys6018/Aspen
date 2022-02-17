package main

type Environment struct {
	enclosing *Environment
	values    map[string]interface{}
}

func (e Environment) Define(name string, value interface{}) {
	e.values[name] = value
}

func (e Environment) Assign(name string, value interface{}) {
	_, ok := e.values[name]

	if ok {
		e.values[name] = value
	} else if e.enclosing != nil {
		e.enclosing.Assign(name, value)
	} else {
		Unreachable("Environment::Assign")
	}
}

func (e Environment) IsDefined(name string) bool {
	_, ok := e.values[name]
	if !ok && e.enclosing != nil {
		return e.enclosing.IsDefined(name)
	}
	return ok
}

func (e Environment) Get(name string) interface{} {
	val, ok := e.values[name]
	if ok {
		return val
	} else if e.enclosing != nil {
		return e.enclosing.Get(name)
	}

	Unreachable("Environment::Get")
	return nil
}

func NewEnvironment(enclosing *Environment) Environment {
	return Environment{values: make(map[string]interface{}), enclosing: enclosing}
}
