#!/bin/bash
##如果task进程已存在，先关闭
net=$(ps aux|grep tasktask)

##去掉空格，并把空格替换为下划线
net=${net// /_}

##以换行符分隔字符串，并把结果赋值给数组arr
OLD_IFS="$IFS"
IFS=$'\n'  # 换行符分隔
arr=($net)
IFS="$OLD_IFS"

##遍历数组arr
for i in ${arr[@]}
do
    len=${#i}
    if [ len > 100 ]
    then
        #删除左数的第一个"__"左边的内容
        c=${i#*"__"}
        #删除右数的最后一个"__"右边的内容
        c=${c%%"__"*}
        c=${c//_/}
        kill $c
        echo "kill $c"
    fi
done    
