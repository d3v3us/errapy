package pkg

import (
	"fmt"
	"runtime"
)

func GetStackTrace(skip int) []string {
	var stackTrace []string
	var pc uintptr
	var line int
	var file string
	var ok bool
	for {
		pc, file, line, ok = runtime.Caller(skip)
		if !ok {
			break
		}
		fn := runtime.FuncForPC(pc)
		funcName := fn.Name()
		stackTrace = append(stackTrace, fmt.Sprintf("%s %s:%d", funcName, file, line))
		skip++
	}
	return stackTrace
}
func FuncForPC(pc uintptr) *runtime.Func {
	return runtime.FuncForPC(pc)
}
func Caller(skip int) (pc uintptr, file string, line int, ok bool) {
	return runtime.Caller(skip)
}
