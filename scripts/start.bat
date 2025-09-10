@echo off
chcp 65001 >nul

echo 启动认证授权中心...

REM 检查Go环境
where go >nul 2>nul
if %errorlevel% neq 0 (
    echo 错误: Go环境未安装
    pause
    exit /b 1
)

REM 检查配置文件
if not exist "config\app.ini" (
    echo 错误: 配置文件 config\app.ini 不存在
    echo 请复制 config\app.ini.example 并修改配置
    pause
    exit /b 1
)

REM 下载依赖
echo 下载依赖...
go mod download

REM 生成Swagger文档
echo 生成API文档...
where swag >nul 2>nul
if %errorlevel% equ 0 (
    swag init -g main.go -o docs
) else (
    echo 警告: swag命令未找到，跳过API文档生成
    echo 安装命令: go install github.com/swaggo/swag/cmd/swag@latest
)

REM 启动服务
echo 启动服务...
go run main.go

pause
