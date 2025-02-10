#!/bin/bash
PROGRAM_NAME=$(ls clash-node-monitor* 2>/dev/null | head -n 1)

# 如果没有找到程序文件，退出脚本
if [ -z "$PROGRAM_NAME" ]; then
    echo "错误：找不到匹配的程序文件"
    exit 1
fi

# 查找进程
pid=$(ps -ef | grep $PROGRAM_NAME | grep -v grep | awk '{print $2}')

if [ -n "$pid" ]; then
    echo "正在停止 $PROGRAM_NAME 进程，PID: $pid"
    kill -15 $pid
    
    # 等待进程结束
    sleep 2
    
    # 检查进程是否还在运行
    if ps -p $pid > /dev/null; then
        echo "进程未能正常停止，正在强制终止..."
        kill -9 $pid
    fi
    
    echo "$PROGRAM_NAME 已停止"
else
    echo "$PROGRAM_NAME 当前未运行"
fi