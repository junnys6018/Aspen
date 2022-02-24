package main

type TypeCastEntry struct {
	from, to *Type
	handler  func(interface{}) interface{}
}

var TypeCasts = make([]TypeCastEntry, 0)

func AddConversion(from, to *Type, handler func(from interface{}) interface{}) {
	TypeCasts = append(TypeCasts, TypeCastEntry{from, to, handler})
}

func GetHandler(from, to *Type) func(from interface{}) interface{} {
	for i := range TypeCasts {
		e := &TypeCasts[i]
		if TypesEqual(from, e.from) && TypesEqual(to, e.to) {
			return e.handler
		}
	}
	Unreachable("type_cast.go: GetHandler")
	return nil
}

func IsConversionLegal(from, to *Type) bool {
	for i := range TypeCasts {
		e := &TypeCasts[i]
		if TypesEqual(from, e.from) && TypesEqual(to, e.to) {
			return true
		}
	}
	return false
}
