# LanLink 状态检查脚本
# 用于快速检查 LanLink 运行状态

Write-Host "`n====== LanLink 状态检查 ======`n" -ForegroundColor Cyan

# 1. 检查进程
Write-Host "1. 进程状态:" -ForegroundColor Yellow
$process = Get-Process lanlink -ErrorAction SilentlyContinue
if ($process) {
    Write-Host "   ✓ 进程运行中 (PID: $($process.Id))" -ForegroundColor Green
    Write-Host "   内存占用: $([math]::Round($process.WorkingSet64/1MB, 2)) MB" -ForegroundColor Gray
} else {
    Write-Host "   ✗ 进程未运行" -ForegroundColor Red
}

# 2. 检查日志
Write-Host "`n2. 日志状态:" -ForegroundColor Yellow
if (Test-Path "lanlink.log") {
    $lastModified = (Get-Item "lanlink.log").LastWriteTime
    $secondsAgo = ((Get-Date) - $lastModified).TotalSeconds
    Write-Host "   ✓ 日志文件存在" -ForegroundColor Green
    Write-Host "   最后更新: $([math]::Round($secondsAgo)) 秒前" -ForegroundColor Gray
    
    if ($secondsAgo -lt 30) {
        Write-Host "   ✓ 日志活跃（系统正常运行）" -ForegroundColor Green
    } else {
        Write-Host "   ⚠ 日志停滞（可能异常）" -ForegroundColor Yellow
    }
    
    Write-Host "`n   最新日志:" -ForegroundColor Gray
    Get-Content "lanlink.log" -Tail 5 | ForEach-Object { 
        Write-Host "   $_" -ForegroundColor Gray 
    }
} else {
    Write-Host "   ✗ 日志文件不存在" -ForegroundColor Red
}

# 3. 检查 Hosts 文件
Write-Host "`n3. Hosts 文件:" -ForegroundColor Yellow
$hostsPath = "C:\Windows\System32\drivers\etc\hosts"
try {
    $hostsContent = Get-Content $hostsPath
    $inManagedZone = $false
    $managedEntries = @()

    foreach ($line in $hostsContent) {
        if ($line -match "LanLink Managed Begin") {
            $inManagedZone = $true
        } elseif ($line -match "LanLink Managed End") {
            $inManagedZone = $false
        } elseif ($inManagedZone -and $line -match "^\d+\.\d+\.\d+\.\d+\s+\S+") {
            $managedEntries += $line.Trim()
        }
    }

    if ($managedEntries.Count -gt 0) {
        Write-Host "   ✓ 发现 $($managedEntries.Count) 个设备:" -ForegroundColor Green
        foreach ($entry in $managedEntries) {
            Write-Host "     $entry" -ForegroundColor Gray
        }
    } else {
        Write-Host "   ⚠ 未发现任何设备（可能还未发现其他节点）" -ForegroundColor Yellow
    }

    # 4. 网络测试
    Write-Host "`n4. 网络测试:" -ForegroundColor Yellow
    if ($managedEntries.Count -gt 0) {
        $testEntry = $managedEntries[0] -split '\s+'
        $testDomain = $testEntry[1]
        Write-Host "   测试 $testDomain ..." -ForegroundColor Gray
        $pingResult = Test-Connection $testDomain -Count 1 -Quiet -ErrorAction SilentlyContinue
        if ($pingResult) {
            Write-Host "   ✓ Ping 成功（域名解析正常）" -ForegroundColor Green
        } else {
            Write-Host "   ⚠ Ping 失败（设备可能离线）" -ForegroundColor Yellow
        }
    } else {
        Write-Host "   ⚠ 无可测试的设备" -ForegroundColor Yellow
    }
} catch {
    Write-Host "   ✗ 无法读取 Hosts 文件（需要管理员权限）" -ForegroundColor Red
}

# 5. 配置检查
Write-Host "`n5. 配置状态:" -ForegroundColor Yellow
if (Test-Path "config.json") {
    Write-Host "   ✓ 配置文件存在" -ForegroundColor Green
    try {
        $config = Get-Content "config.json" | ConvertFrom-Json
        Write-Host "   设备名: $($config.deviceName)" -ForegroundColor Gray
        Write-Host "   域名后缀: $($config.domainSuffix)" -ForegroundColor Gray
        Write-Host "   心跳间隔: $($config.heartbeatIntervalSec) 秒" -ForegroundColor Gray
    } catch {
        Write-Host "   ⚠ 配置文件格式错误" -ForegroundColor Yellow
    }
} else {
    Write-Host "   ⚠ 使用默认配置" -ForegroundColor Yellow
}

# 总结
Write-Host "`n=============================" -ForegroundColor Cyan
Write-Host "状态总结:" -ForegroundColor Yellow

$allGood = $true
if (-not $process) {
    Write-Host "  ✗ 程序未运行" -ForegroundColor Red
    $allGood = $false
}
if (-not (Test-Path "lanlink.log")) {
    Write-Host "  ✗ 日志文件不存在" -ForegroundColor Red
    $allGood = $false
}

if ($allGood) {
    Write-Host "  ✓ 系统运行正常" -ForegroundColor Green
} else {
    Write-Host "  ⚠ 检测到问题，请检查上述输出" -ForegroundColor Yellow
}

Write-Host "`n提示: 使用 'Get-Content lanlink.log -Wait' 实时查看日志" -ForegroundColor Gray
Write-Host "=============================" -ForegroundColor Cyan
Write-Host ""

