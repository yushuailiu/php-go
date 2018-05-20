#include <php.h>
#ifndef PHPGO_ZVAL_H
#define PHPGO_ZVAL_H

void *get_zval_str_val(zval *zval_value);

int get_zval_int_val(zval *zval_value);

double get_zval_double_val(zval *zval_value);


int get_zval_type(zval *zval_value);
#endif