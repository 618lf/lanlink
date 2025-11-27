@echo off
REM LanLink 编译脚本 (Windows)

echo 开始编译 LanLink...

REM 创建输出目录
if not exist dist mkdir dist

REM 编译当前平台
echo 编译 Windows (amd64)...
go build -o dist\lanlink.exe

REM 编译 Linux 版本
echo 编译 Linux (amd64)...
set GOOS=linux
set GOARCH=amd64
go build -o dist\lanlink-linux

REM 编译 macOS 版本
echo 编译 macOS (arm64)...
set GOOS=darwin
set GOARCH=arm64
go build -o dist\lanlink-mac-arm64

echo 编译 macOS (amd64)...
set GOOS=darwin
set GOARCH=amd64
go build -o dist\lanlink-mac-amd64

REM 重置环境变量
set GOOS=
set GOARCH=

echo 编译完成！输出目录: dist\
dir dist\

pause

