package lang


/*

#include "extension.h"

*/
import "C"
import (
	"unsafe"
)

type Extension struct {
	name string
	version string
	constants []*Constant
	functions map[string]*Function
	classes map[string]*Class
}

var extension *Extension

func NewExtension(name string, version string) *Extension {

	extension = &Extension{name:name,version:version}
	return extension
}

func (extension *Extension)GetModule() unsafe.Pointer  {
	initialFunctions()

	registerFunctions()

	me := C.get_zend_module_entry(C.CString(extension.name), C.CString(extension.version))
	return unsafe.Pointer(me)
}

func initialFunctions() {
	C.initialFunctions(funcPointer(module_startup_func), funcPointer(module_shutdown_func),
		funcPointer(request_startup_func), funcPointer(request_shutdown_func))
}

func module_startup_func(moduleType int, moduleNumber int) int {
	println("module_startup_func1")

	registerConstants(moduleType, moduleNumber)

	return 0
}



func module_shutdown_func(moduleType int, moduleNumber int) int {
	println("module_shutdown_func")
	return 0
}
func request_startup_func(moduleType int, moduleNumber int) int {
	println("request_startup_func")
	return 0
}

func request_shutdown_func(moduleType int, moduleNumber int) int {
	println("request_shutdown_func")
	return 0
}

func funcPointer(f interface{}) unsafe.Pointer {
	return unsafe.Pointer(&f)
}