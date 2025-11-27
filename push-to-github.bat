@echo off
chcp 65001 >nul
echo ========================================
echo   LanLink 推送到 GitHub
echo ========================================
echo.

echo 当前远程仓库:
git remote -v
echo.

echo 当前分支和状态:
git status
echo.

echo ========================================
echo 开始推送到 GitHub...
echo ========================================
echo.

git push -u origin main

if %ERRORLEVEL% EQU 0 (
    echo.
    echo ✅ 推送成功！
    echo.
    echo 访问: https://github.com/618lf/lanlink
    echo.
) else (
    echo.
    echo ❌ 推送失败！
    echo.
    echo 可能的原因:
    echo 1. 网络问题 - 检查网络连接或使用代理
    echo 2. 认证问题 - 配置 GitHub 凭据
    echo 3. 超时问题 - 文件过大或网络不稳定
    echo.
    echo 解决方案请查看: 推送指南.txt
    echo.
)

pause

