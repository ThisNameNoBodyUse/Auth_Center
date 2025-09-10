#!/bin/bash

# 认证授权中心启动脚本

echo "启动认证授权中心..."

# 检查Go环境
if ! command -v go &> /dev/null; then
    echo "错误: Go环境未安装"
    exit 1
fi

# 检查Go版本
GO_VERSION=$(go version | cut -d' ' -f3 | sed 's/go//')
REQUIRED_VERSION="1.24.2"

if [ "$(printf '%s\n' "$REQUIRED_VERSION" "$GO_VERSION" | sort -V | head -n1)" != "$REQUIRED_VERSION" ]; then
    echo "错误: Go版本需要 $REQUIRED_VERSION 或更高版本，当前版本: $GO_VERSION"
    exit 1
fi

# 检查配置文件
if [ ! -f "config/app.ini" ]; then
    echo "错误: 配置文件 config/app.ini 不存在"
    echo "请复制 config/app.ini.example 并修改配置"
    exit 1
fi

# 检查数据库连接
echo "检查数据库连接..."
# 这里可以添加数据库连接检查逻辑

# 检查Redis连接
echo "检查Redis连接..."
# 这里可以添加Redis连接检查逻辑

# 下载依赖
echo "下载依赖..."
go mod download

# 生成Swagger文档
echo "生成API文档..."
if command -v swag &> /dev/null; then
    swag init -g main.go -o docs
else
    echo "警告: swag命令未找到，跳过API文档生成"
    echo "安装命令: go install github.com/swaggo/swag/cmd/swag@latest"
fi

# 启动服务
echo "启动服务..."
go run main.go
