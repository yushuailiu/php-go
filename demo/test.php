<?php
function echoLine($msg){
    echo $msg . PHP_EOL;
}


echoLine("常量测试");
echoLine("string constant:" . CONSTANT1);
echoLine("stringl constant:" . CONSTANT2);
echoLine("int constant:" . INTCONSTANT);
echoLine("bool constant:" . var_export(BOOLCONSTANT, true));
echoLine("float constant:" . FLOATCONSTANT);
echoLine("null constant:" . var_export(NULLCONSTANT, true));


echoLine("方法测试");
helloWorld();
