@echo off
chcp 65001 >nul
title LanLink 全平台编译

echo.
echo ============================================
echo   LanLink 全平台编译
echo ============================================
echo.

REM 创建输出目录
if not exist dist mkdir dist

REM 清理旧文件
del /q dist\* 2>nul

echo [1/4] 编译 Windows (amd64)...
set GOOS=windows
set GOARCH=amd64
go build -ldflags="-s -w" -o dist\lanlink-windows-amd64.exe
if %errorlevel% equ 0 (
    echo   √ Windows 编译成功
) else (
    echo   × Windows 编译失败
)

echo.
echo [2/4] 编译 Linux (amd64)...
set GOOS=linux
set GOARCH=amd64
go build -ldflags="-s -w" -o dist\lanlink-linux-amd64
if %errorlevel% equ 0 (
    echo   √ Linux 编译成功
) else (
    echo   × Linux 编译失败
)

echo.
echo [3/4] 编译 Mac Intel (amd64)...
set GOOS=darwin
set GOARCH=amd64
go build -ldflags="-s -w" -o dist\lanlink-mac-amd64
if %errorlevel% equ 0 (
    echo   √ Mac Intel 编译成功
) else (
    echo   × Mac Intel 编译失败
)

echo.
echo [4/4] 编译 Mac Apple Silicon (arm64)...
set GOOS=darwin
set GOARCH=arm64
go build -ldflags="-s -w" -o dist\lanlink-mac-arm64
if %errorlevel% equ 0 (
    echo   √ Mac Apple Silicon 编译成功
) else (
    echo   × Mac Apple Silicon 编译失败
)

REM 重置环境变量
set GOOS=
set GOARCH=

echo.
echo ============================================
echo   编译完成！
echo ============================================
echo.
echo 输出目录: dist\
echo.
dir dist\ /b
echo.
echo 平台说明:
echo   lanlink-windows-amd64.exe  - Windows (所有PC)
echo   lanlink-linux-amd64        - Linux (大部分PC)
echo   lanlink-mac-amd64          - Mac Intel (2020年前)
echo   lanlink-mac-arm64          - Mac M1/M2/M3
echo.
echo 使用方法:
echo   1. 将对应平台的文件传输到目标设备
echo   2. 在目标设备上运行: sudo ./lanlink install
echo   3. 设置开机自启: sudo lanlink service install
echo.
pause

