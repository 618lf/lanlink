# LanLink 全平台编译脚本
# 在 Windows 上一键编译所有平台版本

# 创建输出目录
New-Item -ItemType Directory -Force -Path dist | Out-Null

Write-Host "`n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
Write-Host "  LanLink 全平台编译" -ForegroundColor Cyan
Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan

# Windows amd64
Write-Host "`n[1/4] 编译 Windows (amd64)..." -ForegroundColor Yellow
$env:GOOS="windows"; $env:GOARCH="amd64"
go build -ldflags="-s -w" -o dist/lanlink-windows-amd64.exe
if ($LASTEXITCODE -eq 0) { 
    Write-Host "  ✓ Windows 编译成功" -ForegroundColor Green 
} else {
    Write-Host "  ✗ Windows 编译失败" -ForegroundColor Red
}

# Linux amd64
Write-Host "`n[2/4] 编译 Linux (amd64)..." -ForegroundColor Yellow
$env:GOOS="linux"; $env:GOARCH="amd64"
go build -ldflags="-s -w" -o dist/lanlink-linux-amd64
if ($LASTEXITCODE -eq 0) { 
    Write-Host "  ✓ Linux 编译成功" -ForegroundColor Green 
} else {
    Write-Host "  ✗ Linux 编译失败" -ForegroundColor Red
}

# Mac Intel (amd64)
Write-Host "`n[3/4] 编译 Mac Intel (amd64)..." -ForegroundColor Yellow
$env:GOOS="darwin"; $env:GOARCH="amd64"
go build -ldflags="-s -w" -o dist/lanlink-mac-amd64
if ($LASTEXITCODE -eq 0) { 
    Write-Host "  ✓ Mac Intel 编译成功" -ForegroundColor Green 
} else {
    Write-Host "  ✗ Mac Intel 编译失败" -ForegroundColor Red
}

# Mac Apple Silicon (arm64)
Write-Host "`n[4/4] 编译 Mac Apple Silicon (arm64)..." -ForegroundColor Yellow
$env:GOOS="darwin"; $env:GOARCH="arm64"
go build -ldflags="-s -w" -o dist/lanlink-mac-arm64
if ($LASTEXITCODE -eq 0) { 
    Write-Host "  ✓ Mac Apple Silicon 编译成功" -ForegroundColor Green 
} else {
    Write-Host "  ✗ Mac Apple Silicon 编译失败" -ForegroundColor Red
}

# 重置环境变量
$env:GOOS=""; $env:GOARCH=""

Write-Host "`n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
Write-Host "  编译完成！" -ForegroundColor Green
Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan

Write-Host "`nOutput directory: dist/`n" -ForegroundColor Yellow
Get-ChildItem dist/ | Select-Object Name, @{
    Name="Size(MB)"; Expression={[math]::Round($_.Length / 1MB, 2)}
}, LastWriteTime | Format-Table -AutoSize

Write-Host "Platform Guide:" -ForegroundColor Cyan
Write-Host "  Windows PC        -> lanlink-windows-amd64.exe" -ForegroundColor Gray
Write-Host "  Linux PC          -> lanlink-linux-amd64" -ForegroundColor Gray
Write-Host "  Mac Intel         -> lanlink-mac-amd64" -ForegroundColor Gray
Write-Host "  Mac M1/M2/M3      -> lanlink-mac-arm64" -ForegroundColor Gray
Write-Host ""

