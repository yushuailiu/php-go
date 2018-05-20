package lang


/*
#include <zend.h>
#include "extension.h"

#include "zval.h"

#include <zend_types.h>

 */
import "C"
import (
	"unsafe"
	"reflect"
	"log"
)

type Function struct {
	Name string
	Handler interface{}
}


func (extension *Extension) RegisterFunction(function *Function) {
	class, ok := extension.classes["global"];
	if !ok {
		class = &Class{Name: "global"}
		extension.classes = make(map[string]*Class)
		extension.classes["global"] = class
	}
	_, ok = class.Methods[function.Name]
	if ok {
		panic("repeat register function:" + function.Name)
	}
	if class.Methods == nil {
		class.Methods = make(map[string]*Function);
	}
	class.Methods[function.Name] = function
}


func registerFunctions() {
	class, ok := extension.classes["global"]
	if !ok {
		return
	}

	for _, function := range class.Methods {
		function.register()
	}
}

func (function *Function) register() {
	C.zend_add_function(C.CString(function.Name))
}

//export functionParamNum
func functionParamNum(functionName *C.char) C.int {
	functionStringName := C.GoString(functionName)
	fn := getFunctionByName(functionStringName)
	funcType := reflect.TypeOf((fn.Handler))
	return C.int(funcType.NumIn())
}

//export functionRetType
func functionRetType(functionName *C.char) C.int {
	functionStringName := C.GoString(functionName)
	fn := getFunctionByName(functionStringName)
	funcType := reflect.TypeOf(fn.Handler)
	if reflect.TypeOf(fn.Handler).NumOut() == 0{
		return 0
	}
	goType := funcType.Out(0).Kind()
	cType := goType2PhpType(goType)
	return C.int(cType)
}

func goType2PhpType(goType interface{}) C.int {

	switch goType {
	case reflect.Invalid:
		return C.IS_NULL
	case reflect.String:
		return C.IS_STRING
	case reflect.Bool:
		return C._IS_BOOL
	case reflect.Int64:
		fallthrough
	case reflect.Uint64:
		fallthrough
	case reflect.Int32:
		fallthrough
	case reflect.Uint32:
		fallthrough
	case reflect.Int:
		fallthrough
	case reflect.Uint:
		fallthrough
	case reflect.Int16:
		fallthrough
	case reflect.Uint16:
		fallthrough
	case reflect.Int8:
		fallthrough
	case reflect.Uint8:
		return C.IS_LONG
	case reflect.Float32, reflect.Float64:
		return C.IS_DOUBLE
	case reflect.Slice:
		return C.IS_ARRAY
	case reflect.Map:
		return C.IS_ARRAY
	}
	return C.int(0)
}

//export callGoFunction
func callGoFunction(functionName *C.char, args []C.zval) uintptr {
	functionStringName := C.GoString(functionName)
	fn := getFunctionByName(functionStringName)
	funcType := reflect.TypeOf(fn.Handler)
	argTypes := make([]reflect.Value, 0)
	for index := 0; index < funcType.NumIn(); index++ {
		switch funcType.In(index).Kind() {
		case reflect.String:
			if C.get_zval_type(&args[index]) != C.IS_STRING {
				panic("parameter 1 must be zend_array.")
			}
			argTypes = append(argTypes, reflect.ValueOf(C.GoString((*C.char)(C.get_zval_str_val(&args[index])))))
		case reflect.Float64:
			fallthrough
		case reflect.Float32:
			argTypes = append(argTypes, reflect.ValueOf(C.get_zval_double_val(&args[index])).Convert(funcType.In(index)))
		case reflect.Bool:
			switch reflect.ValueOf(C.get_zval_type(&args[index])).Int() {
			case C.IS_FALSE:
				argTypes = append(argTypes, reflect.ValueOf(false))
			case C.IS_TRUE:
				argTypes = append(argTypes, reflect.ValueOf(true))
			}
		case reflect.Int64:
			fallthrough
		case reflect.Uint64:
			fallthrough
		case reflect.Int32:
			fallthrough
		case reflect.Uint32:
			fallthrough
		case reflect.Int:
			fallthrough
		case reflect.Uint:
			fallthrough
		case reflect.Int16:
			fallthrough
		case reflect.Uint16:
			fallthrough
		case reflect.Int8:
			fallthrough
		case reflect.Uint8:
			argTypes = append(argTypes, reflect.ValueOf(C.get_zval_int_val(&args[index])).Convert(funcType.In(index)))
		case reflect.Ptr:
		case reflect.Interface:
		case reflect.Slice:
		case reflect.Map:
		default:
			log.Panicln("Unsupported kind:", "wtf", funcType.In(index).String(),
				funcType.In(index).Kind(), reflect.TypeOf(fn.Handler).IsVariadic())
		}
	}

	result := reflect.ValueOf(fn.Handler).Call(argTypes)

	ret :=  transformRet2PhpVal(fn.Handler, result)
	return ret
}

func transformRet2PhpVal(fn interface{}, result []reflect.Value) (retp uintptr) {

	funcType := reflect.TypeOf(fn);

	if funcType.Kind() != reflect.Func {
		panic("the type of fn must be func")
	}

	if funcType.NumOut() > 0 {

		switch result[0].Kind() {
		case reflect.String:
			retp = uintptr(unsafe.Pointer(C.CString(result[0].Interface().(string))))
		case reflect.Float64:
			fallthrough
		case reflect.Float32:
			var pdv *C.double = (*C.double)(C.malloc(8))
			*pdv = (C.double)(result[0].Interface().(float64))
			retp = uintptr(unsafe.Pointer(pdv))
		case reflect.Bool:
			var bv = result[0].Interface().(bool)
			if bv {
				retp = uintptr(1)
			} else {
				retp = uintptr(0)
			}
		case reflect.Int64:
			fallthrough
		case reflect.Uint64:
			fallthrough
		case reflect.Int32:
			fallthrough
		case reflect.Uint32:
			fallthrough
		case reflect.Int:
			fallthrough
		case reflect.Uint:
			fallthrough
		case reflect.Int16:
			fallthrough
		case reflect.Uint16:
			fallthrough
		case reflect.Int8:
			fallthrough
		case reflect.Uint8:
			var ret *C.int = (*C.int)(C.malloc(4))
			*ret = (C.int)(result[0].Int())
			retp = uintptr(unsafe.Pointer(ret))
		case reflect.Ptr:
			var nv = result[0].Pointer()
			retp = uintptr(nv)
		default:
			log.Panicln("unsupport return type:", funcType.Out(0).Kind().String())

		}

	}

	return retp
}

func getFunctionByName(functionName string) (fn *Function){
	class, ok := extension.classes["global"]
	if !ok {
		panic("no function has register")
	}
	fn, ok = class.Methods[functionName]
	if !ok {
		panic("function is not register:" + functionName)
	}
	return fn
}

//export functionParamInfo
func functionParamInfo(fp unsafe.Pointer, class string) *C.char {
	funcType := reflect.TypeOf(*(*interface{})(fp))
	argTypes := ""
	for index := 0; index < funcType.NumIn(); index++ {
		switch funcType.In(index).Kind() {
		case reflect.String:
			argTypes = argTypes + "s"
		case reflect.Float64:
			fallthrough
		case reflect.Float32:
			argTypes = argTypes + "d"
		case reflect.Bool:
			argTypes = argTypes + "b"
		case reflect.Int64:
			fallthrough
		case reflect.Uint64:
			fallthrough
		case reflect.Int32:
			fallthrough
		case reflect.Uint32:
			fallthrough
		case reflect.Int:
			fallthrough
		case reflect.Uint:
			fallthrough
		case reflect.Int16:
			fallthrough
		case reflect.Uint16:
			fallthrough
		case reflect.Int8:
			fallthrough
		case reflect.Uint8:
			argTypes = argTypes + "l"
		case reflect.Ptr:
			argTypes = argTypes + "p"
		case reflect.Interface:
			argTypes = argTypes + "a" // Any/interface
		case reflect.Slice:
			argTypes = argTypes + "v" // vector
		case reflect.Map:
			argTypes = argTypes + "m"
		default:
			log.Panicln("Unsupported kind:", "wtf", funcType.In(index).String(),
				funcType.In(index).Kind(), reflect.TypeOf(*(*interface{})(fp)).IsVariadic())
		}
	}
	return C.CString(argTypes)
}