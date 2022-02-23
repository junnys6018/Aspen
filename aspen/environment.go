package main

type Environment struct {
	enclosing *Environment
	values    map[string]interface{}
}

func NewGlobalEnvironment() Environment {
	environment := NewEnvironment(nil)

	// copy native functions into global environment
	for k, v := range NativeFunctions {
		environment.values[k] = v
	}

	return environment
}

func (e Environment) Define(name string, value interface{}) {
	e.values[name] = value
}

func (e Environment) Ancestor(depth int) Environment {
	environment := e
	for i := 0; i < depth; i++ {
		if environment.enclosing == nil {
			panic("Environment::Ancestor: bad depth")
		}
		environment = *environment.enclosing
	}

	return environment
}

func (e Environment) AssignAt(name string, depth int, value interface{}) {
	ancestor := e.Ancestor(depth)

	if !ancestor.IsDefinedLocally(name) {
		panic("Environment::AssignAt: name not defined")
	}

	ancestor.values[name] = value
}

func (e Environment) IsDefined(name string) bool {
	_, ok := e.values[name]
	if !ok && e.enclosing != nil {
		return e.enclosing.IsDefined(name)
	}
	return ok
}

func (e Environment) IsDefinedLocally(name string) bool {
	_, ok := e.values[name]
	return ok
}

func (e Environment) GetAt(name string, depth int) interface{} {
	return e.Ancestor(depth).values[name]
}

func (e Environment) GetDepth(name string) int {
	if e.IsDefinedLocally(name) {
		return 0
	} else if e.enclosing != nil {
		return 1 + e.enclosing.GetDepth(name)
	}

	Unreachable("Environment::GetDepth")
	return 0
}

func NewEnvironment(enclosing *Environment) Environment {
	return Environment{values: make(map[string]interface{}), enclosing: enclosing}
}
