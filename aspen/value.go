package main

import "fmt"

func IsValueType(iface interface{}) bool {
	switch iface.(type) {
	case int64, uint64, float64, bool, nil:
		return true
	default:
		return false
	}
}

func PrintValue(iface interface{}) {
	switch v := iface.(type) {
	case []rune:
		fmt.Println(string(v))
	default:
		fmt.Println(v)
	}
}

func ValuesEqual(lhs, rhs interface{}) bool {
	// todo: revisit this when we implement object types
	if IsValueType(lhs) {
		return lhs == rhs
	}

	switch lhsV := lhs.(type) {
	case []rune:
		rhsV := rhs.([]rune)
		if len(rhsV) != len(lhsV) {
			return false
		}

		for i := range lhsV {
			if lhsV[i] != rhsV[i] {
				return false
			}
		}
		return true
	}

	Unreachable("ValuesEqual")
	return false
}
func OperatorGreater(lhs, rhs interface{}) interface{} {
	switch lhsV := lhs.(type) {
	case int64:
		rhsV := rhs.(int64)
		return lhsV > rhsV
	case uint64:
		rhsV := rhs.(uint64)
		return lhsV > rhsV
	case float64:
		rhsV := rhs.(float64)
		return lhsV > rhsV
	}
	Unreachable("value.go: OperatorGreater")
	return nil
}
func OperatorGreaterEqual(lhs, rhs interface{}) interface{} {
	switch lhsV := lhs.(type) {
	case int64:
		rhsV := rhs.(int64)
		return lhsV >= rhsV
	case uint64:
		rhsV := rhs.(uint64)
		return lhsV >= rhsV
	case float64:
		rhsV := rhs.(float64)
		return lhsV >= rhsV
	}
	Unreachable("value.go: OperatorGreaterEqual")
	return nil
}
func OperatorLess(lhs, rhs interface{}) interface{} {
	switch lhsV := lhs.(type) {
	case int64:
		rhsV := rhs.(int64)
		return lhsV < rhsV
	case uint64:
		rhsV := rhs.(uint64)
		return lhsV < rhsV
	case float64:
		rhsV := rhs.(float64)
		return lhsV < rhsV
	}
	Unreachable("value.go: OperatorLess")
	return nil
}
func OperatorLessEqual(lhs, rhs interface{}) interface{} {
	switch lhsV := lhs.(type) {
	case int64:
		rhsV := rhs.(int64)
		return lhsV <= rhsV
	case uint64:
		rhsV := rhs.(uint64)
		return lhsV <= rhsV
	case float64:
		rhsV := rhs.(float64)
		return lhsV <= rhsV
	}
	Unreachable("value.go: OperatorLessEqual")
	return nil
}
func OperatorPipe(lhs, rhs interface{}) interface{} {
	switch lhsV := lhs.(type) {
	case int64:
		rhsV := rhs.(int64)
		return lhsV | rhsV
	case uint64:
		rhsV := rhs.(uint64)
		return lhsV | rhsV
	}
	Unreachable("value.go: OperatorPipe")
	return nil
}
func OperatorCaret(lhs, rhs interface{}) interface{} {
	switch lhsV := lhs.(type) {
	case int64:
		rhsV := rhs.(int64)
		return lhsV ^ rhsV
	case uint64:
		rhsV := rhs.(uint64)
		return lhsV ^ rhsV
	}
	Unreachable("value.go: OperatorCaret")
	return nil
}
func OperatorAmp(lhs, rhs interface{}) interface{} {
	switch lhsV := lhs.(type) {
	case int64:
		rhsV := rhs.(int64)
		return lhsV & rhsV
	case uint64:
		rhsV := rhs.(uint64)
		return lhsV & rhsV
	}
	Unreachable("value.go: OperatorAmp")
	return nil
}
func OperatorMinus(lhs, rhs interface{}) interface{} {
	switch lhsV := lhs.(type) {
	case int64:
		rhsV := rhs.(int64)
		return lhsV - rhsV
	case uint64:
		rhsV := rhs.(uint64)
		return lhsV - rhsV
	case float64:
		rhsV := rhs.(float64)
		return lhsV - rhsV
	}
	Unreachable("value.go: OperatorMinus")
	return nil
}
func OperatorPlus(lhs, rhs interface{}) interface{} {
	switch lhsV := lhs.(type) {
	case int64:
		rhsV := rhs.(int64)
		return lhsV + rhsV
	case uint64:
		rhsV := rhs.(uint64)
		return lhsV + rhsV
	case float64:
		rhsV := rhs.(float64)
		return lhsV + rhsV
	}
	Unreachable("value.go: OperatorPlus")
	return nil
}
func OperatorSlash(lhs, rhs interface{}) interface{} {
	switch lhsV := lhs.(type) {
	case int64:
		rhsV := rhs.(int64)
		return lhsV / rhsV
	case uint64:
		rhsV := rhs.(uint64)
		return lhsV / rhsV
	case float64:
		rhsV := rhs.(float64)
		return lhsV / rhsV
	}
	Unreachable("value.go: OperatorSlash")
	return nil
}
func OperatorStar(lhs, rhs interface{}) interface{} {
	switch lhsV := lhs.(type) {
	case int64:
		rhsV := rhs.(int64)
		return lhsV * rhsV
	case uint64:
		rhsV := rhs.(uint64)
		return lhsV * rhsV
	case float64:
		rhsV := rhs.(float64)
		return lhsV * rhsV
	}
	Unreachable("value.go: OperatorStar")
	return nil
}
func OperatorModulus(lhs, rhs interface{}) interface{} {
	switch lhsV := lhs.(type) {
	case int64:
		rhsV := rhs.(int64)
		return lhsV % rhsV
	case uint64:
		rhsV := rhs.(uint64)
		return lhsV % rhsV
	}
	Unreachable("value.go: OperatorModulus")
	return nil
}
