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



void phpgo_function_handler(zend_execute_data *execute_data, zval *return_value)
{
    printf("asdfasdfasdf");
}

void zend_add_function(char *name)
{
    zend_function_entry e = {"helloWorld", phpgo_function_handler, NULL, 0, 0};
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


