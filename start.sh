#!/bin/bash

PROGRAM_NAME="clash-node-monitor"

# 查找并杀掉已有的进程
pid=$(ps -ef | grep $PROGRAM_NAME | grep -v grep | awk '{print $2}')
if [ -n "$pid" ]; then
    echo "正在终止已运行的 $PROGRAM_NAME 进程，PID: $pid"
    kill -9 $pid
else
    echo "$PROGRAM_NAME 当前未运行"
fi

# 检查程序文件是否存在
if [ ! -f "./$PROGRAM_NAME" ]; then
    echo "错误：找不到 $PROGRAM_NAME 程序文件"
    exit 1
fi

# 检查配置文件是否存在
if [ ! -f "./config.yaml" ]; then
    echo "错误：找不到 config.yaml 配置文件"
    exit 1
fi

# 设置可执行权限
chmod +x ./$PROGRAM_NAME

# 循环启动程序，如果程序崩溃则自动重启
while true
do
    echo "正在启动 $PROGRAM_NAME..."
    ./$PROGRAM_NAME &
    if [ $? -ne 0 ]; then
        echo "$PROGRAM_NAME 运行异常，退出代码 $?，1秒后尝试重启..." >&2
        sleep 1
    else
        break
    fi
done

echo "$PROGRAM_NAME 已成功启动"