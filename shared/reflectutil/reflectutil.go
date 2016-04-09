package reflectutil

import (
	"reflect"
	"runtime"
)

// FuncName returns the function name for the given func.
func FuncName(f interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}
