#include <php.h>

#ifdef ZTS
#include "TSRM.h"
#endif

// add this head file then can call go functions
#include "_cgo_export.h"

static int(*go_module_startup_func)(int, int) = 0;
static int(*go_module_shutdown_func)(int, int) = 0;
static int(*go_request_startup_func)(int, int) = 0;
static int(*go_request_shutdown_func)(int, int) = 0;

static zend_class_entry g_entry = {0};

zend_function_entry *functions;
int functions_num = 0;

void go_function_ret_transform(int go_return_type, void *retp, zval *return_value) {
    switch (go_return_type) {
        case IS_STRING:
            RETVAL_STRINGL((char*)retp, strlen((char*)retp));
            break;
        case IS_DOUBLE:
            RETVAL_DOUBLE(*(double*)retp);
            free((double*)retp);
            break;
        case IS_TRUE:
            RETVAL_TRUE;
            break;
        case IS_FALSE:
            RETVAL_FALSE;
            break;
        case _IS_BOOL:
            RETVAL_BOOL(*(int*)retp);
            break;
        case IS_LONG:
            RETVAL_LONG(*(int*)retp);
            break;
        case IS_NULL:
            RETVAL_NULL();
            break;
        case IS_UNDEF:
            RETVAL_NULL();
            break;
        case IS_RESOURCE:
            RETVAL_NULL();
            break;
        case IS_ARRAY:
            break;
        default:
            zend_error(E_WARNING, "unknown type of return value: %d .", go_return_type);
            break;
        }
}

void phpgo_function_handler(zend_execute_data *execute_data, zval *return_value)
{
    zend_string *function_name_temp = execute_data->func->common.function_name;
    char *function_name = ZSTR_VAL(function_name_temp);
    char *class_name = ZEND_FN_SCOPE_NAME(execute_data->func);
    int num_args = functionParamNum(function_name);
    GoSlice args_go;
    args_go.data = NULL;
    args_go.len = num_args;
    args_go.cap = num_args;

    zval *args;
    if (num_args) {
        args = (zval*)malloc( num_args*sizeof(zval) );
        args_go.data = args;
    }

    if (zend_get_parameters_array_ex(num_args, args) == FAILURE) {
        WRONG_PARAM_COUNT;
        return;
    }

    void *retp = (void*)callGoFunction(function_name, args_go);
    int go_return_type = functionRetType(function_name);

    if (go_return_type == 0) {
        return;
    }

    go_function_ret_transform(go_return_type, retp, return_value);
}


void zend_add_function(char *name)
{
    zend_function_entry e = {name, phpgo_function_handler, NULL, 0, 0};
    functions_num ++;
    if (functions_num == 1) {
        functions = (zend_function_entry*)calloc(functions_num + 1, sizeof(zend_function_entry));
    } else {
        functions = (zend_function_entry*)realloc(functions, (functions_num + 1) * sizeof(zend_function_entry));
    }
    memset(&functions[functions_num], 0, sizeof(zend_function_entry));
    memcpy(&functions[functions_num-1], &e, sizeof(zend_function_entry));
}


void initialFunctions(void *module_startup_func, void *module_shutdown_func, void *request_startup_func, void *request_shutdown_func)
{
    go_module_startup_func = (int(*)(int, int))module_startup_func;
    go_module_shutdown_func = (int(*)(int, int))module_shutdown_func;
    go_request_startup_func = (int(*)(int, int))request_startup_func;
    go_request_shutdown_func = (int(*)(int, int))request_shutdown_func;
}

int module_startup_func(int type, int module_number)
{
    if (go_module_startup_func) {
        return call_go_func(go_module_startup_func, type, module_number);
    }
    return 0;
}

int module_shutdown_func(int type, int module_number)
{
    if (go_module_shutdown_func) {
        return call_go_func(go_module_shutdown_func, type, module_number);
    }
    return 0;
}
int request_startup_func(int type, int module_number)
{
    if (go_request_startup_func) {
        return call_go_func(go_request_startup_func, type, module_number);
    }
    return 0;
}
int request_shutdown_func(int type, int module_number)
{
    if (go_request_shutdown_func) {
        return call_go_func(go_request_shutdown_func, type, module_number);
    }
    return 0;
}


zend_module_entry ze;

zend_module_entry *get_zend_module_entry(char *name, char *version) {

zend_module_entry te = {
        STANDARD_MODULE_HEADER,
        name,
        functions,
        module_startup_func,
        module_shutdown_func,
        request_startup_func,
        request_shutdown_func,
        NULL,
        version,
        STANDARD_MODULE_PROPERTIES
    };
    memcpy(&ze, &te, sizeof(zend_module_entry));
    return &ze;
}


