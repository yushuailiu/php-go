<?php
echo "扩展注册的函数列表:" . PHP_EOL;
$ext_info = array();
$loaded_extensions=get_loaded_extensions();//获取已加载的扩展
foreach($loaded_extensions as $ext)
{
    if ($ext !== "demo") {
        continue;
    }
    $funs=get_extension_funcs($ext);//获取某一扩展下的所有函数
    if(!empty($funs) && is_array($funs))
    {
        foreach($funs as $fun)
        {
            $reflect = new ReflectionFunction($fun);
            $params = $reflect->getParameters();//获取函数参数信息
            $param_str = '';
            if(!empty($params) && is_array($params))
            {
                foreach($params as $param) {
                    if($param->getName() != '')
                    {
                        $param_str .= '$'.$param->getName().',';
                    }
                }
                $param_str = substr($param_str,0,-1);
            }
            echo $fun.'('.$param_str.')' . PHP_EOL;
        }
    }
}

echo "\n";
echo "\n";
