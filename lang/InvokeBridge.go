package lang

import (
	"unsafe"
	"reflect"
	"log"
)


/*

#include <stdlib.h>

*/
import "C"

//export call_go_func
func call_go_func(fp unsafe.Pointer, moduleType uintptr, moduleNumber uintptr) uintptr {
	callFunc := reflect.ValueOf(*(*interface{})(fp))

	args := []uintptr{moduleType, moduleNumber}
	argv := transformPhpArgsToGoArgs(callFunc.Interface(), args)
	callFunc.Call(argv)
	return 0
}


func transformPhpArgsToGoArgs(fn interface{}, args []uintptr) (argv []reflect.Value) {
	fty := reflect.TypeOf(fn)
	if fty.Kind() != reflect.Func {
		log.Panicln("the pointer call by c not a func pointer:", fty.Kind().String())
	}

	argv = make([]reflect.Value, 0)
	for idx := 0; idx < fty.NumIn(); idx++ {
		switch fty.In(idx).Kind() {
		case reflect.String:
			var arg = C.GoString((*C.char)(unsafe.Pointer(args[idx])))
			var v = reflect.ValueOf(arg).Convert(fty.In(idx))
			argv = append(argv, v)
		case reflect.Float64, reflect.Float32:
			var arg = (C.double)(args[idx])
			var v = reflect.ValueOf(arg).Convert(fty.In(idx))
			argv = append(argv, v)
			C.free(unsafe.Pointer(args[idx]))
		case reflect.Bool:
			var arg = (C.int)(args[idx])
			if arg == 1 {
				argv = append(argv, reflect.ValueOf(true))
			} else {
				argv = append(argv, reflect.ValueOf(false))
			}
		case reflect.Int64, reflect.Uint64:
			fallthrough
		case reflect.Int32, reflect.Uint32:
			fallthrough
		case reflect.Int, reflect.Uint:
			fallthrough
		case reflect.Int16, reflect.Uint16:
			fallthrough
		case reflect.Int8, reflect.Uint8:
			var arg = (C.longlong)(args[idx])
			var v = reflect.ValueOf(arg).Convert(fty.In(idx))
			argv = append(argv, v)
		default:
			log.Panicln("field type not support:",
				fty.In(idx).Kind().String(), fty.In(idx).String())
		}
	}

	if len(argv) != fty.NumIn() {
		panic("arguments number is wrong")
	}
	return
}