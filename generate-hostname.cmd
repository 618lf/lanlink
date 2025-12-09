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

REM 获取 BIOS 序列号（使用 /value 格式避免编码问题）
for /f "tokens=2 delims==" %%a in ('wmic bios get serialnumber /value 2^>nul ^| find "="') do (
    set "SERIAL=%%a"
)
REM 去除回车符
set "SERIAL=!SERIAL: =!"
for /f "delims=" %%a in ("!SERIAL!") do set "SERIAL=%%a"

REM 检查是否为无效值
if /i "!SERIAL!"=="To Be Filled By O.E.M." set "SERIAL="
if /i "!SERIAL!"=="Default string" set "SERIAL="
if /i "!SERIAL!"=="None" set "SERIAL="
if "!SERIAL!"=="" goto :try_baseboard
goto :process

:try_baseboard
REM 尝试主板序列号
for /f "tokens=2 delims==" %%a in ('wmic baseboard get serialnumber /value 2^>nul ^| find "="') do (
    set "SERIAL=%%a"
)
set "SERIAL=!SERIAL: =!"
for /f "delims=" %%a in ("!SERIAL!") do set "SERIAL=%%a"

if /i "!SERIAL!"=="To Be Filled By O.E.M." set "SERIAL="
if /i "!SERIAL!"=="Default string" set "SERIAL="
if /i "!SERIAL!"=="None" set "SERIAL="
if "!SERIAL!"=="" goto :try_uuid
goto :process

:try_uuid
REM 使用 UUID
for /f "tokens=2 delims==" %%a in ('wmic csproduct get uuid /value 2^>nul ^| find "="') do (
    set "SERIAL=%%a"
)
set "SERIAL=!SERIAL: =!"
set "SERIAL=!SERIAL:-=!"
for /f "delims=" %%a in ("!SERIAL!") do set "SERIAL=%%a"

:process
REM 转小写
set "CLEAN=!SERIAL!"
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

REM 取后6位（使用负数索引）
set "HOSTID=!CLEAN:~-6!"

REM 如果为空或不足6位，用0填充
if "!HOSTID!"=="" set "HOSTID=000000"
:pad_loop
if "!HOSTID:~5,1!"=="" (
    set "HOSTID=0!HOSTID!"
    goto :pad_loop
)

echo win-!HOSTID!.%DOMAIN_SUFFIX%

endlocal
