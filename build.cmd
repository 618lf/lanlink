@echo off
chcp 65001 >nul
title LanLink 全平台编译

echo.
echo ============================================
echo   LanLink 全平台编译
echo ============================================
echo.

REM 清理并创建输出目录
if exist dist rmdir /s /q dist
mkdir dist\windows-amd64
mkdir dist\linux-amd64
mkdir dist\mac-amd64
mkdir dist\mac-arm64

echo [1/4] 编译 Windows (amd64)...
set GOOS=windows
set GOARCH=amd64
go build -ldflags="-s -w" -o dist\windows-amd64\lanlink.exe
if %errorlevel% equ 0 (
    echo   √ Windows 编译成功
) else (
    echo   × Windows 编译失败
)

echo.
echo [2/4] 编译 Linux (amd64)...
set GOOS=linux
set GOARCH=amd64
go build -ldflags="-s -w" -o dist\linux-amd64\lanlink
if %errorlevel% equ 0 (
    echo   √ Linux 编译成功
) else (
    echo   × Linux 编译失败
)

echo.
echo [3/4] 编译 Mac Intel (amd64)...
set GOOS=darwin
set GOARCH=amd64
go build -ldflags="-s -w" -o dist\mac-amd64\lanlink
if %errorlevel% equ 0 (
    echo   √ Mac Intel 编译成功
) else (
    echo   × Mac Intel 编译失败
)

echo.
echo [4/4] 编译 Mac Apple Silicon (arm64)...
set GOOS=darwin
set GOARCH=arm64
go build -ldflags="-s -w" -o dist\mac-arm64\lanlink
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
echo 输出目录结构:
echo   dist\
echo     windows-amd64\lanlink.exe
echo     linux-amd64\lanlink
echo     mac-amd64\lanlink
echo     mac-arm64\lanlink
echo.
echo 平台说明:
echo   windows-amd64  - Windows (所有PC)
echo   linux-amd64    - Linux (大部分PC)
echo   mac-amd64      - Mac Intel (2020年前)
echo   mac-arm64      - Mac M1/M2/M3/M4
echo.
pause
