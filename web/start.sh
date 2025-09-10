#!/bin/bash

echo "正在启动认证授权中心管理后台..."
echo

# 检查 Node.js 是否安装
if ! command -v node &> /dev/null; then
    echo "错误: 未检测到 Node.js，请先安装 Node.js"
    echo "下载地址: https://nodejs.org/"
    exit 1
fi

# 检查 npm 是否安装
if ! command -v npm &> /dev/null; then
    echo "错误: 未检测到 npm，请检查 Node.js 安装"
    exit 1
fi

# 检查是否已安装依赖
if [ ! -d "node_modules" ]; then
    echo "正在安装依赖包..."
    npm install
    if [ $? -ne 0 ]; then
        echo "错误: 依赖安装失败"
        exit 1
    fi
    echo "依赖安装完成"
    echo
fi

echo "启动开发服务器..."
echo "访问地址: http://localhost:3000"
echo "按 Ctrl+C 停止服务器"
echo

npm run dev
