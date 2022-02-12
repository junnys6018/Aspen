package main

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
