@echo off
chcp 65001 >nul
title LanLink 推送到 GitHub

echo.
echo ╔══════════════════════════════════════════╗
echo ║     LanLink 推送到 GitHub 工具           ║
echo ╚══════════════════════════════════════════╝
echo.

:menu
echo 请选择推送方式:
echo.
echo [1] 直接推送 (无代理)
echo [2] 使用 HTTP 代理推送
echo [3] 查看当前状态
echo [4] 配置 GitHub Token 后推送
echo [5] 退出
echo.
set /p choice=请输入选项 (1-5): 

if "%choice%"=="1" goto direct
if "%choice%"=="2" goto proxy
if "%choice%"=="3" goto status
if "%choice%"=="4" goto token
if "%choice%"=="5" goto end
goto menu

:direct
echo.
echo ► 正在直接推送...
git push -u origin main
if %ERRORLEVEL% EQU 0 (
    echo.
    echo ✅ 推送成功！
    echo 访问: https://github.com/618lf/lanlink
    goto success
) else (
    echo.
    echo ❌ 推送失败！请尝试其他方式
    goto menu
)

:proxy
echo.
echo 请输入代理地址 (例如: http://127.0.0.1:7890)
set /p proxy_url=代理地址: 
echo.
echo ► 配置代理: %proxy_url%
git config --global http.proxy %proxy_url%
git config --global https.proxy %proxy_url%
echo.
echo ► 正在推送...
git push -u origin main
set push_result=%ERRORLEVEL%
echo.
echo ► 清除代理配置...
git config --global --unset http.proxy
git config --global --unset https.proxy
if %push_result% EQU 0 (
    echo.
    echo ✅ 推送成功！
    echo 访问: https://github.com/618lf/lanlink
    goto success
) else (
    echo.
    echo ❌ 推送失败！请检查代理设置
    goto menu
)

:status
echo.
echo ► 当前 Git 状态:
echo.
git status
echo.
echo ► 远程仓库:
git remote -v
echo.
echo ► 本地提交:
git log --oneline -3
echo.
pause
goto menu

:token
echo.
echo ► 配置 GitHub Personal Access Token
echo.
echo 步骤:
echo 1. 访问: https://github.com/settings/tokens
echo 2. 点击 "Generate new token (classic)"
echo 3. 勾选 "repo" 权限
echo 4. 生成并复制 token
echo.
set /p token=请粘贴您的 token: 
echo.
echo ► 配置凭据...
git remote set-url origin https://%token%@github.com/618lf/lanlink.git
echo.
echo ► 正在推送...
git push -u origin main
if %ERRORLEVEL% EQU 0 (
    echo.
    echo ✅ 推送成功！
    echo 访问: https://github.com/618lf/lanlink
    goto success
) else (
    echo.
    echo ❌ 推送失败！请检查 token 是否正确
    git remote set-url origin https://github.com/618lf/lanlink.git
    goto menu
)

:success
echo.
echo ════════════════════════════════════════
echo   推送成功！接下来你可以:
echo ════════════════════════════════════════
echo.
echo 1. 访问仓库: https://github.com/618lf/lanlink
echo 2. 发布版本:
echo    git tag -a v1.0.0 -m "Release v1.0.0"
echo    git push origin v1.0.0
echo.
pause
goto end

:end
echo.
echo 感谢使用！
exit

