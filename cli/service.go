package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// ServiceInstall 安装为系统服务（开机自启）
func ServiceInstall() error {
	if runtime.GOOS == "windows" {
		return serviceInstallWindows()
	}
	return serviceInstallUnix()
}

// ServiceUninstall 卸载系统服务
func ServiceUninstall() error {
	if runtime.GOOS == "windows" {
		return serviceUninstallWindows()
	}
	return serviceUninstallUnix()
}

// ServiceStart 启动服务
func ServiceStart() error {
	if runtime.GOOS == "windows" {
		return serviceStartWindows()
	}
	return serviceStartUnix()
}

// ServiceStop 停止服务
func ServiceStop() error {
	if runtime.GOOS == "windows" {
		return serviceStopWindows()
	}
	return serviceStopUnix()
}

// ServiceStatus 查看服务状态
func ServiceStatus() error {
	if runtime.GOOS == "windows" {
		return serviceStatusWindows()
	}
	return serviceStatusUnix()
}

// ===== Windows 实现 =====

func serviceInstallWindows() error {
	Header("安装 LanLink 服务 (Windows)")

	// 检查管理员权限
	Section("检查权限")
	if !isAdmin() {
		Error("需要管理员权限")
		fmt.Println("\n请以管理员身份运行此命令")
		return fmt.Errorf("需要管理员权限")
	}
	Success("已获取管理员权限")

	// 获取程序路径
	exePath, err := os.Executable()
	if err != nil {
		Error("获取程序路径失败: %v", err)
		return err
	}
	exePath, _ = filepath.Abs(exePath)

	Section("创建服务")
	KeyValue("程序路径", exePath)

	// 使用 sc 命令创建 Windows 服务
	cmd := exec.Command("sc", "create", "LanLink",
		"binPath=", fmt.Sprintf("\"%s\" start", exePath),
		"start=", "auto",
		"DisplayName=", "LanLink - 局域网域名自动映射")

	output, err := cmd.CombinedOutput()
	if err != nil {
		Error("创建服务失败: %v", err)
		fmt.Println(string(output))
		return err
	}
	Success("服务创建成功")

	// 设置服务描述
	cmd = exec.Command("sc", "description", "LanLink",
		"局域网域名自动映射工具，自动发现和同步局域网设备")
	cmd.Run()

	Footer()

	fmt.Println()
	Success("服务安装完成！")
	fmt.Println("\n启动服务:")
	fmt.Println("  lanlink service start")
	fmt.Println("  或: sc start LanLink")
	fmt.Println("\n服务将在开机时自动启动")
	fmt.Println()

	return nil
}

func serviceUninstallWindows() error {
	Header("卸载 LanLink 服务 (Windows)")

	// 检查管理员权限
	Section("检查权限")
	if !isAdmin() {
		Error("需要管理员权限")
		return fmt.Errorf("需要管理员权限")
	}
	Success("已获取管理员权限")

	// 停止服务
	Section("停止服务")
	cmd := exec.Command("sc", "stop", "LanLink")
	if err := cmd.Run(); err != nil {
		Warn("停止服务失败（可能未运行）")
	} else {
		Success("服务已停止")
	}

	// 删除服务
	Section("删除服务")
	cmd = exec.Command("sc", "delete", "LanLink")
	if err := cmd.Run(); err != nil {
		Error("删除服务失败: %v", err)
		return err
	}
	Success("服务已删除")

	Footer()

	fmt.Println()
	Success("服务卸载完成！")
	fmt.Println()

	return nil
}

func serviceStartWindows() error {
	Info("启动 LanLink 服务...")

	cmd := exec.Command("sc", "start", "LanLink")
	output, err := cmd.CombinedOutput()
	if err != nil {
		Error("启动失败: %v", err)
		fmt.Println(string(output))
		return err
	}

	Success("服务已启动")
	fmt.Println("\n查看状态: lanlink service status")
	return nil
}

func serviceStopWindows() error {
	Info("停止 LanLink 服务...")

	cmd := exec.Command("sc", "stop", "LanLink")
	output, err := cmd.CombinedOutput()
	if err != nil {
		Error("停止失败: %v", err)
		fmt.Println(string(output))
		return err
	}

	Success("服务已停止")
	return nil
}

func serviceStatusWindows() error {
	Header("LanLink 服务状态 (Windows)")

	cmd := exec.Command("sc", "query", "LanLink")
	output, err := cmd.CombinedOutput()

	if err != nil {
		Warn("服务未安装")
		fmt.Println("\n安装服务: lanlink service install")
		return nil
	}

	fmt.Println(string(output))
	Footer()

	return nil
}

// ===== Unix (Linux/Mac) 实现 =====

func serviceInstallUnix() error {
	Header("安装 LanLink 服务 (Systemd)")

	// 检查 root 权限
	Section("检查权限")
	if os.Geteuid() != 0 {
		Error("需要 root 权限")
		fmt.Println("\n请使用 sudo 运行:")
		fmt.Println("  sudo lanlink service install")
		return fmt.Errorf("需要 root 权限")
	}
	Success("已获取 root 权限")

	// 获取程序路径
	exePath, err := os.Executable()
	if err != nil {
		Error("获取程序路径失败: %v", err)
		return err
	}
	exePath, _ = filepath.Abs(exePath)

	Section("创建服务文件")
	KeyValue("程序路径", exePath)

	// 创建 systemd 服务文件
	serviceContent := fmt.Sprintf(`[Unit]
Description=LanLink - 局域网域名自动映射工具
After=network.target

[Service]
Type=simple
ExecStart=%s start
Restart=on-failure
RestartSec=10
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
`, exePath)

	servicePath := "/etc/systemd/system/lanlink.service"
	if err := os.WriteFile(servicePath, []byte(serviceContent), 0644); err != nil {
		Error("创建服务文件失败: %v", err)
		return err
	}
	Success("服务文件创建成功")
	KeyValue("服务文件", servicePath)

	// 重载 systemd
	Section("配置服务")
	cmd := exec.Command("systemctl", "daemon-reload")
	if err := cmd.Run(); err != nil {
		Warn("重载 systemd 失败: %v", err)
	} else {
		Success("Systemd 已重载")
	}

	// 启用服务
	cmd = exec.Command("systemctl", "enable", "lanlink")
	if err := cmd.Run(); err != nil {
		Error("启用服务失败: %v", err)
		return err
	}
	Success("服务已启用（开机自启）")

	Footer()

	fmt.Println()
	Success("服务安装完成！")
	fmt.Println("\n启动服务:")
	fmt.Println("  lanlink service start")
	fmt.Println("  或: sudo systemctl start lanlink")
	fmt.Println("\n查看日志:")
	fmt.Println("  sudo journalctl -u lanlink -f")
	fmt.Println("\n服务将在开机时自动启动")
	fmt.Println()

	return nil
}

func serviceUninstallUnix() error {
	Header("卸载 LanLink 服务 (Systemd)")

	// 检查 root 权限
	Section("检查权限")
	if os.Geteuid() != 0 {
		Error("需要 root 权限")
		fmt.Println("\n请使用 sudo 运行:")
		fmt.Println("  sudo lanlink service uninstall")
		return fmt.Errorf("需要 root 权限")
	}
	Success("已获取 root 权限")

	// 停止服务
	Section("停止服务")
	cmd := exec.Command("systemctl", "stop", "lanlink")
	if err := cmd.Run(); err != nil {
		Warn("停止服务失败（可能未运行）")
	} else {
		Success("服务已停止")
	}

	// 禁用服务
	Section("禁用服务")
	cmd = exec.Command("systemctl", "disable", "lanlink")
	if err := cmd.Run(); err != nil {
		Warn("禁用服务失败（可能未启用）")
	} else {
		Success("服务已禁用")
	}

	// 删除服务文件
	Section("删除服务文件")
	servicePath := "/etc/systemd/system/lanlink.service"
	if err := os.Remove(servicePath); err != nil {
		Warn("删除服务文件失败: %v", err)
	} else {
		Success("服务文件已删除")
	}

	// 重载 systemd
	cmd = exec.Command("systemctl", "daemon-reload")
	cmd.Run()

	Footer()

	fmt.Println()
	Success("服务卸载完成！")
	fmt.Println()

	return nil
}

func serviceStartUnix() error {
	Info("启动 LanLink 服务...")

	cmd := exec.Command("systemctl", "start", "lanlink")
	output, err := cmd.CombinedOutput()
	if err != nil {
		Error("启动失败: %v", err)
		fmt.Println(string(output))
		fmt.Println("\n提示: 需要 root 权限，请使用 sudo")
		return err
	}

	Success("服务已启动")
	fmt.Println("\n查看状态: lanlink service status")
	fmt.Println("查看日志: sudo journalctl -u lanlink -f")
	return nil
}

func serviceStopUnix() error {
	Info("停止 LanLink 服务...")

	cmd := exec.Command("systemctl", "stop", "lanlink")
	output, err := cmd.CombinedOutput()
	if err != nil {
		Error("停止失败: %v", err)
		fmt.Println(string(output))
		fmt.Println("\n提示: 需要 root 权限，请使用 sudo")
		return err
	}

	Success("服务已停止")
	return nil
}

func serviceStatusUnix() error {
	Header("LanLink 服务状态 (Systemd)")

	cmd := exec.Command("systemctl", "status", "lanlink", "--no-pager")
	output, err := cmd.CombinedOutput()

	if err != nil {
		// 即使服务未运行，status 也可能返回错误，但我们仍然显示输出
		fmt.Println(string(output))
		if len(output) == 0 {
			Warn("服务未安装")
			fmt.Println("\n安装服务: sudo lanlink service install")
		}
	} else {
		fmt.Println(string(output))
	}

	Footer()

	return nil
}

