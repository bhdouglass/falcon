package scopes

// #include "shim.h"
import "C"

// These functions are used by tests.  They are not part of a
// *_test.go file because they make use of cgo.

func newTestingResult() *Result {
	return makeResult(C.new_testing_result())
}
