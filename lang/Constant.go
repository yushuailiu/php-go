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
}

func (extension *Extension) RegisterConstant(name string, val interface{}, namespace interface{}) {
	constant := &Constant{name, val, namespace}
	extension.constants = append(extension.constants, constant)
}

func (constant *Constant) register(moduleType int, moduleNumber int) {

	switch reflect.ValueOf(constant.Val).Kind() {
	case reflect.String:
		zendRegisterStringConstant(moduleNumber, constant)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32,
		reflect.Uint64:
			
	case reflect.Bool:
	case reflect.Float32, reflect.Float64:
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
	C.zend_register_string_constant(cname, C.size_t(len(constant.Name)), C.CString(val), C.CONST_CS|C.CONST_PERSISTENT, C.int(moduleNumber));
}
