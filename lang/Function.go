package lang


/*
#include "extension.h"
 */
import "C"

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
	C.zend_add_function(C.CString("test"))
}