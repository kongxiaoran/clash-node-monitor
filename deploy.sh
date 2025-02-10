#!/bin/bash

# 配置信息
REMOTE_USER="root"
REMOTE_HOST="K6ymycZlK500!23jKS"
REMOTE_DIR="/app/clash-node-monitor"
BINARY_NAME="clash-node-monitor"

# 编译程序
echo "开始编译程序..."
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${BINARY_NAME} .

if [ $? -ne 0 ]; then
    echo "编译失败！"
    exit 1
fi

# 创建远程目录前先清理
echo "清理远程目录..."
ssh ${REMOTE_USER}@${REMOTE_HOST} "if [ -d ${REMOTE_DIR} ]; then rm -rf ${REMOTE_DIR}/*; else mkdir -p ${REMOTE_DIR}; fi"

# 上传文件
echo "上传文件到服务器..."
scp ${BINARY_NAME} Klink.yaml config.yaml start.sh ${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_DIR}/

if [ $? -ne 0 ]; then
    echo "文件上传失败！"
    exit 1
fi

# 设置权限
echo "设置执行权限..."
ssh ${REMOTE_USER}@${REMOTE_HOST} "chmod +x ${REMOTE_DIR}/${BINARY_NAME} ${REMOTE_DIR}/start.sh"

# 重启服务
echo "重启服务..."
ssh ${REMOTE_USER}@${REMOTE_HOST} "cd ${REMOTE_DIR} && ./start.sh"

echo "部署完成！"