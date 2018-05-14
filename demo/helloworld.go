package main

import "C"

import (
	"unsafe"
	"../lang"
)
//export get_module
func get_module() unsafe.Pointer {
	extension := lang.NewExtension("test", "0.0.1")

	extension.RegisterConstant("helloworld", "hello world!!", nil)
	return extension.GetModule()
}
func main() {  }