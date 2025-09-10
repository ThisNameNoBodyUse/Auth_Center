@echo off
echo 正在启动认证授权中心管理后台...
echo.

REM 检查 Node.js 是否安装
node --version >nul 2>&1
if %errorlevel% neq 0 (
    echo 错误: 未检测到 Node.js，请先安装 Node.js
    echo 下载地址: https://nodejs.org/
    pause
    exit /b 1
)

REM 检查 npm 是否安装
npm --version >nul 2>&1
if %errorlevel% neq 0 (
    echo 错误: 未检测到 npm，请检查 Node.js 安装
    pause
    exit /b 1
)

REM 检查是否已安装依赖
if not exist "node_modules" (
    echo 正在安装依赖包...
    npm install
    if %errorlevel% neq 0 (
        echo 错误: 依赖安装失败
        pause
        exit /b 1
    )
    echo 依赖安装完成
    echo.
)

echo 启动开发服务器...
echo 访问地址: http://localhost:3000
echo 按 Ctrl+C 停止服务器
echo.

npm run dev
