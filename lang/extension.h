
#include <php.h>

#ifndef PHPGO_EXTENSION_H
#define PHPGO_EXTENSION_H

extern void* phpgo_get_module(char *name, char *version);
extern int phpgo_get_module_number();
void phpgo_register_init_functions(void *module_startup_func, void *module_shutdown_func,
                                   void *request_startup_func, void *request_shutdown_func);
extern int zend_add_class(int cidx, char *name);
extern int zend_add_method(int cidx, int fidx, int cbid, char *name, char *mname, char *atys, int rety);
void zend_add_function(char *name);
void initialFunctions(void *module_startup_func, void *module_shutdown_func, void *request_startup_func, void *request_shutdown_func);

zend_module_entry *get_zend_module_entry(char *, char *);

extern zend_function_entry *functions;

typedef int (*intFunc) (int, int);

#endif
