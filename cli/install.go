package cli

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// Install 安装 LanLink 到系统
func Install() error {
	if runtime.GOOS == "windows" {
		return installWindows()
	}
	return installUnix()
}

// installWindows Windows 安装
func installWindows() error {
	Header("安装 LanLink (Windows)")

	// 1. 检查管理员权限
	Section("检查权限")
	if !isAdmin() {
		Error("需要管理员权限")
		fmt.Println("\n请以管理员身份运行此命令:")
		fmt.Println("  1. 右键点击 PowerShell")
		fmt.Println("  2. 选择\"以管理员身份运行\"")
		fmt.Println("  3. 运行: lanlink install")
		return fmt.Errorf("需要管理员权限")
	}
	Success("已获取管理员权限")

	// 2. 获取当前程序路径
	exePath, err := os.Executable()
	if err != nil {
		Error("获取程序路径失败: %v", err)
		return err
	}
	exePath, _ = filepath.Abs(exePath)

	Section("程序信息")
	KeyValue("当前路径", exePath)

	// 3. 创建安装目录
	installDir := `C:\Program Files\LanLink`
	Section("安装目录")
	KeyValue("目标位置", installDir)

	if err := os.MkdirAll(installDir, 0755); err != nil {
		Error("创建目录失败: %v", err)
		return err
	}
	Success("目录创建成功")

	// 4. 复制程序文件
	Section("复制文件")
	targetPath := filepath.Join(installDir, "lanlink.exe")

	if err := copyFile(exePath, targetPath); err != nil {
		Error("复制文件失败: %v", err)
		return err
	}
	Success("程序文件已复制")

	// 5. 添加到系统 PATH
	Section("配置环境变量")
	if err := addToPath(installDir); err != nil {
		Error("添加到PATH失败: %v", err)
		return err
	}
	Success("已添加到系统 PATH")

	// 6. 创建启动脚本（可选）
	Section("创建快捷方式")
	if err := createStartupScript(installDir); err != nil {
		Warn("创建启动脚本失败: %v", err)
	} else {
		Success("已创建启动脚本")
	}

	Footer()

	// 安装完成提示
	fmt.Println()
	Success("安装完成!")
	fmt.Println("\n现在您可以在任何位置使用 'lanlink' 命令了")
	fmt.Println("\n重启终端或运行以下命令刷新环境变量:")
	fmt.Println("  $env:Path = [System.Environment]::GetEnvironmentVariable(\"Path\",\"Machine\")")
	fmt.Println("\n测试安装:")
	fmt.Println("  lanlink version")
	fmt.Println()

	return nil
}

// installUnix Linux/Mac 安装
func installUnix() error {
	Header("安装 LanLink (Unix)")

	// 1. 检查 root 权限
	Section("检查权限")
	if os.Geteuid() != 0 {
		Error("需要 root 权限")
		fmt.Println("\n请使用 sudo 运行:")
		fmt.Println("  sudo lanlink install")
		return fmt.Errorf("需要 root 权限")
	}
	Success("已获取 root 权限")

	// 2. 获取当前程序路径
	exePath, err := os.Executable()
	if err != nil {
		Error("获取程序路径失败: %v", err)
		return err
	}
	exePath, _ = filepath.Abs(exePath)

	Section("程序信息")
	KeyValue("当前路径", exePath)

	// 3. 复制到 /usr/local/bin
	targetPath := "/usr/local/bin/lanlink"
	Section("安装")
	KeyValue("目标位置", targetPath)

	if err := copyFile(exePath, targetPath); err != nil {
		Error("复制文件失败: %v", err)
		return err
	}

	// 4. 设置执行权限
	if err := os.Chmod(targetPath, 0755); err != nil {
		Error("设置权限失败: %v", err)
		return err
	}
	Success("安装完成")

	// 5. 创建 systemd 服务（可选）
	Section("系统服务")
	if err := createSystemdService(); err != nil {
		Warn("创建系统服务失败: %v", err)
		fmt.Println("  您仍然可以手动运行: sudo lanlink")
	} else {
		Success("已创建 systemd 服务")
		fmt.Println("\n启动服务:")
		fmt.Println("  sudo systemctl start lanlink")
		fmt.Println("  sudo systemctl enable lanlink  # 开机自启")
	}

	Footer()

	// 安装完成提示
	fmt.Println()
	Success("安装完成!")
	fmt.Println("\n现在您可以在任何位置使用 'lanlink' 命令了")
	fmt.Println("\n测试安装:")
	fmt.Println("  lanlink version")
	fmt.Println()

	return nil
}

// copyFile 复制文件
func copyFile(src, dst string) error {
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	return err
}

// addToPath 添加到系统 PATH (Windows)
func addToPath(dir string) error {
	// 使用 PowerShell 添加到系统 PATH
	script := fmt.Sprintf(`
$path = [Environment]::GetEnvironmentVariable("Path", "Machine")
if ($path -notlike "*%s*") {
    $newPath = $path + ";%s"
    [Environment]::SetEnvironmentVariable("Path", $newPath, "Machine")
    Write-Host "PATH updated"
} else {
    Write-Host "Already in PATH"
}
`, dir, dir)

	cmd := exec.Command("powershell", "-Command", script)
	return cmd.Run()
}

// createStartupScript 创建启动脚本 (Windows)
func createStartupScript(installDir string) error {
	script := fmt.Sprintf(`@echo off
REM LanLink 启动脚本
echo 正在启动 LanLink...
cd /d "%s"
lanlink.exe
pause
`, installDir)

	scriptPath := filepath.Join(installDir, "start-lanlink.bat")
	return os.WriteFile(scriptPath, []byte(script), 0644)
}

// createSystemdService 创建 systemd 服务 (Linux)
func createSystemdService() error {
	service := `[Unit]
Description=LanLink - 局域网域名自动映射工具
After=network.target

[Service]
Type=simple
ExecStart=/usr/local/bin/lanlink
Restart=on-failure
RestartSec=10

[Install]
WantedBy=multi-user.target
`

	servicePath := "/etc/systemd/system/lanlink.service"
	if err := os.WriteFile(servicePath, []byte(service), 0644); err != nil {
		return err
	}

	// 重载 systemd
	cmd := exec.Command("systemctl", "daemon-reload")
	return cmd.Run()
}

// isAdmin 检查是否有管理员权限 (Windows)
func isAdmin() bool {
	if runtime.GOOS != "windows" {
		return os.Geteuid() == 0
	}

	// Windows: 尝试执行需要管理员权限的操作
	cmd := exec.Command("net", "session")
	err := cmd.Run()
	return err == nil
}

