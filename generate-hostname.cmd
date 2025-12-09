@echo off
REM
REM generate-hostname.cmd - 基于主板序列号生成唯一主机域名 (Windows CMD)
REM
REM 使用方式:
REM   generate-hostname.cmd    输出域名
REM
REM 域名格式: win-{序列号后6位}.coobee.local
REM

setlocal enabledelayedexpansion

set "DOMAIN_SUFFIX=coobee.local"
set "SERIAL="

REM 获取 BIOS 序列号
for /f "skip=1 tokens=*" %%a in ('wmic bios get serialnumber 2^>nul') do (
    if not defined SERIAL (
        set "SERIAL=%%a"
        set "SERIAL=!SERIAL: =!"
    )
)

REM 如果 BIOS 序列号无效，尝试主板序列号
if "!SERIAL!"=="" (
    for /f "skip=1 tokens=*" %%a in ('wmic baseboard get serialnumber 2^>nul') do (
        if not defined SERIAL (
            set "SERIAL=%%a"
            set "SERIAL=!SERIAL: =!"
        )
    )
)

REM 如果还是无效，使用 UUID
if "!SERIAL!"=="" (
    for /f "skip=1 tokens=*" %%a in ('wmic csproduct get uuid 2^>nul') do (
        if not defined SERIAL (
            set "SERIAL=%%a"
            set "SERIAL=!SERIAL: =!"
            set "SERIAL=!SERIAL:-=!"
        )
    )
)

REM 移除特殊字符，只保留字母数字
set "CLEAN="
set "CHARS=abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
for /l %%i in (0,1,99) do (
    set "CHAR=!SERIAL:~%%i,1!"
    if "!CHAR!"=="" goto :done_clean
    echo !CHARS! | findstr /c:"!CHAR!" >nul 2>&1 && set "CLEAN=!CLEAN!!CHAR!"
)
:done_clean

REM 转小写
set "LOWER="
for %%c in (a b c d e f g h i j k l m n o p q r s t u v w x y z 0 1 2 3 4 5 6 7 8 9) do (
    set "CLEAN=!CLEAN:%%c=%%c!"
)
set "CLEAN=!CLEAN:A=a!"
set "CLEAN=!CLEAN:B=b!"
set "CLEAN=!CLEAN:C=c!"
set "CLEAN=!CLEAN:D=d!"
set "CLEAN=!CLEAN:E=e!"
set "CLEAN=!CLEAN:F=f!"
set "CLEAN=!CLEAN:G=g!"
set "CLEAN=!CLEAN:H=h!"
set "CLEAN=!CLEAN:I=i!"
set "CLEAN=!CLEAN:J=j!"
set "CLEAN=!CLEAN:K=k!"
set "CLEAN=!CLEAN:L=l!"
set "CLEAN=!CLEAN:M=m!"
set "CLEAN=!CLEAN:N=n!"
set "CLEAN=!CLEAN:O=o!"
set "CLEAN=!CLEAN:P=p!"
set "CLEAN=!CLEAN:Q=q!"
set "CLEAN=!CLEAN:R=r!"
set "CLEAN=!CLEAN:S=s!"
set "CLEAN=!CLEAN:T=t!"
set "CLEAN=!CLEAN:U=u!"
set "CLEAN=!CLEAN:V=v!"
set "CLEAN=!CLEAN:W=w!"
set "CLEAN=!CLEAN:X=x!"
set "CLEAN=!CLEAN:Y=y!"
set "CLEAN=!CLEAN:Z=z!"

REM 取后6位
set "LEN=0"
for /l %%i in (0,1,99) do (
    set "CHAR=!CLEAN:~%%i,1!"
    if not "!CHAR!"=="" set /a "LEN=%%i+1"
)

if !LEN! gtr 6 (
    set /a "START=!LEN!-6"
    set "HOSTID=!CLEAN:~%START%,6!"
) else (
    set "HOSTID=!CLEAN!"
)

REM 不足6位用0填充
:pad_loop
set "PADLEN=0"
for /l %%i in (0,1,5) do (
    set "CHAR=!HOSTID:~%%i,1!"
    if not "!CHAR!"=="" set /a "PADLEN=%%i+1"
)
if !PADLEN! lss 6 (
    set "HOSTID=0!HOSTID!"
    goto :pad_loop
)

echo win-!HOSTID!.%DOMAIN_SUFFIX%

endlocal

