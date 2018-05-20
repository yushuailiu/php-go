#include <php.h>

void *get_zval_str_val(zval *zval_value)
{
   return Z_STRVAL(*zval_value);
}

int get_zval_int_val(zval *zval_value)
{
    return Z_LVAL(*zval_value);
}

double get_zval_double_val(zval *zval_value)
{
    return Z_DVAL(*zval_value);
}


int get_zval_type(zval *zval_value)
{
    return Z_TYPE(*zval_value);
}