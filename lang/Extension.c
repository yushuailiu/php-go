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
        NULL,
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


