package main

import "C"

import (
	"unsafe"
	"../lang"
	"fmt"
)
//export get_module
func get_module() unsafe.Pointer {
	extension := lang.NewExtension("test", "0.0.1")

	constant1 := &lang.Constant{Name:"constant1", Val: "hello world!"}
	extension.RegisterConstant(constant1)

	constant2 := &lang.Constant{Name:"constant2", Val: "hello world!", Len:5}
	extension.RegisterConstant(constant2)


	intConstant := &lang.Constant{Name: "INTCONSTANT", Val:6265431}
	extension.RegisterConstant(intConstant)

	boolConstant := &lang.Constant{Name: "BOOLCONSTANT", Val: true}
	extension.RegisterConstant(boolConstant)

	floatConstant := &lang.Constant{Name: "FLOATCONSTANT", Val:123.99}
	extension.RegisterConstant(floatConstant)

	nullConstant := &lang.Constant{Name:"NULLCONSTANT"}
	extension.RegisterConstant(nullConstant)

	helloWorldFunction := &lang.Function{Name:"helloWorld", Handler:helloWorld}
	extension.RegisterFunction(helloWorldFunction)

	helloWorld2Function := &lang.Function{Name:"helloWorld2", Handler:helloWorld2}
	extension.RegisterFunction(helloWorld2Function)

	return extension.GetModule()
}
func main() {  }

func helloWorld(name string, age int, year int, flag bool) int {
	fmt.Printf("hello %s! -- from helloWorld function", name);
	return 25
}

func helloWorld2() string {
	return "I am helloWorld2"
}