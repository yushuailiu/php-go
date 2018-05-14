package lang

/*
#include "extension.h"

#include <php.h>
#include <zend_exceptions.h>
#include <zend_interfaces.h>
#include <zend_ini.h>
#include <zend_constants.h>
#include <SAPI.h>
#include <zend_API.h>

 */
import "C"
import (
	"strings"
	"reflect"
)

type Constant struct {
	Name      string
	Val       interface{}
	Namespace interface{}
	Len int
}

func (extension *Extension) RegisterConstant(constant *Constant) {
	extension.constants = append(extension.constants, constant)
}

func (constant *Constant) register(moduleType int, moduleNumber int) {

	if constant.Val == nil {
		zendRegisterNullConstant(moduleNumber, constant)
	} else {
		switch reflect.ValueOf(constant.Val).Kind() {
		case reflect.String:
			zendRegisterStringConstant(moduleNumber, constant)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32,
			reflect.Uint64:
			zendRegisterIntConstant(moduleNumber, constant)
		case reflect.Bool:
			zendRegisterBoolConstant(moduleNumber, constant)
		case reflect.Float32, reflect.Float64:
			zendRegisterFloatConstant(moduleNumber, constant)
		default:
			panic("not support constant type of " + reflect.TypeOf(constant.Val).Kind().String())
		}
	}

}

func registerConstants(moduleType int, moduleNumber int) {
	for _, constant := range extension.constants {
		constant.register(moduleType, moduleNumber)
	}
}

func zendRegisterStringConstant(moduleNumber int, constant *Constant) {
	cname := C.CString(strings.ToUpper(constant.Name))
	val := constant.Val.(string)
	if constant.Len == 0 {
		C.zend_register_string_constant(cname, C.size_t(len(constant.Name)), C.CString(val), C.CONST_CS|C.CONST_PERSISTENT, C.int(moduleNumber));
	} else {
		C.zend_register_stringl_constant(cname, C.size_t(len(constant.Name)), C.CString(val), C.size_t(constant.Len), C.CONST_CS|C.CONST_PERSISTENT, C.int(moduleNumber));
	}
}


func zendRegisterIntConstant(moduleNumber int, constant *Constant) {
	cname := C.CString(strings.ToUpper(constant.Name))
	val := C.zend_long(reflect.ValueOf(constant.Val).Convert(reflect.TypeOf(int64(1))).Interface().(int64))
	C.zend_register_long_constant(cname, C.size_t(len(constant.Name)), val, C.CONST_CS|C.CONST_PERSISTENT, C.int(moduleNumber))
}

func zendRegisterBoolConstant(moduleNumber int, constant *Constant)  {
	cname := C.CString(strings.ToUpper(constant.Name))
	val := int8(0)
	if constant.Val == true {
		val = 1
	}
	C.zend_register_bool_constant(cname, C.size_t(len(constant.Name)), C.zend_bool(val), C.CONST_CS|C.CONST_PERSISTENT, C.int(moduleNumber))
}

func zendRegisterFloatConstant(moduleNumber int, constant *Constant) {
	cname := C.CString(strings.ToUpper(constant.Name))
	val := C.double(reflect.ValueOf(constant.Val).Convert(reflect.TypeOf(float64(1.0))).Interface().(float64))
	C.zend_register_double_constant(cname, C.size_t(len(constant.Name)), val, C.CONST_CS|C.CONST_PERSISTENT, C.int(moduleNumber))
}

func zendRegisterNullConstant(moduleNumber int, constant *Constant) {
	cname := C.CString(strings.ToUpper(constant.Name))
	C.zend_register_null_constant(cname, C.size_t(len(constant.Name)), C.CONST_CS|C.CONST_PERSISTENT, C.int(moduleNumber))
}