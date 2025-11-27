package cli

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

// Uninstall 卸载 LanLink
func Uninstall() error {
	if runtime.GOOS == "windows" {
		return uninstallWindows()
	}
	return uninstallUnix()
}

// uninstallWindows Windows 卸载
func uninstallWindows() error {
	Header("卸载 LanLink (Windows)")

	// 1. 检查管理员权限
	Section("检查权限")
	if !isAdmin() {
		Error("需要管理员权限")
		fmt.Println("\n请以管理员身份运行此命令")
		return fmt.Errorf("需要管理员权限")
	}
	Success("已获取管理员权限")

	// 2. 从 PATH 移除
	Section("移除环境变量")
	installDir := `C:\Program Files\LanLink`
	if err := removeFromPath(installDir); err != nil {
		Warn("从PATH移除失败: %v", err)
	} else {
		Success("已从系统 PATH 移除")
	}

	// 3. 删除安装目录
	Section("删除文件")
	if err := os.RemoveAll(installDir); err != nil {
		Error("删除目录失败: %v", err)
		return err
	}
	Success("已删除安装目录")

	Footer()

	fmt.Println()
	Success("卸载完成!")
	fmt.Println("\n感谢使用 LanLink")
	fmt.Println()

	return nil
}

// uninstallUnix Linux/Mac 卸载
func uninstallUnix() error {
	Header("卸载 LanLink (Unix)")

	// 1. 检查 root 权限
	Section("检查权限")
	if os.Geteuid() != 0 {
		Error("需要 root 权限")
		fmt.Println("\n请使用 sudo 运行:")
		fmt.Println("  sudo lanlink uninstall")
		return fmt.Errorf("需要 root 权限")
	}
	Success("已获取 root 权限")

	// 2. 停止服务
	Section("停止服务")
	cmd := exec.Command("systemctl", "stop", "lanlink")
	if err := cmd.Run(); err != nil {
		Warn("停止服务失败（可能未运行）")
	} else {
		Success("服务已停止")
	}

	// 3. 禁用服务
	cmd = exec.Command("systemctl", "disable", "lanlink")
	if err := cmd.Run(); err != nil {
		Warn("禁用服务失败（可能未启用）")
	} else {
		Success("服务已禁用")
	}

	// 4. 删除服务文件
	Section("删除文件")
	servicePath := "/etc/systemd/system/lanlink.service"
	if err := os.Remove(servicePath); err != nil {
		Warn("删除服务文件失败: %v", err)
	} else {
		Success("已删除服务文件")
	}

	// 5. 删除程序文件
	targetPath := "/usr/local/bin/lanlink"
	if err := os.Remove(targetPath); err != nil {
		Error("删除程序失败: %v", err)
		return err
	}
	Success("已删除程序文件")

	// 6. 重载 systemd
	cmd = exec.Command("systemctl", "daemon-reload")
	cmd.Run()

	Footer()

	fmt.Println()
	Success("卸载完成!")
	fmt.Println("\n感谢使用 LanLink")
	fmt.Println()

	return nil
}

// removeFromPath 从系统 PATH 移除 (Windows)
func removeFromPath(dir string) error {
	script := fmt.Sprintf(`
$path = [Environment]::GetEnvironmentVariable("Path", "Machine")
$newPath = ($path.Split(';') | Where-Object { $_ -ne '%s' }) -join ';'
[Environment]::SetEnvironmentVariable("Path", $newPath, "Machine")
Write-Host "PATH updated"
`, dir)

	cmd := exec.Command("powershell", "-Command", script)
	return cmd.Run()
}

