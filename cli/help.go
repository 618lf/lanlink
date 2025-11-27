package cli

import (
	"fmt"
)

// ShowHelp 显示帮助信息
func ShowHelp() {
	fmt.Printf(`
%s

用法:
  lanlink [command] [options]

命令:
  start              启动服务（默认）
  status             查看运行状态
  list, ls           列出所有节点
  logs               查看日志
  ping <domain>      测试连接
  install            安装到系统PATH（需要管理员权限）
  uninstall          从系统卸载（需要管理员权限）
  service            服务管理（开机自启）
    install          安装为系统服务
    uninstall        卸载系统服务
    start            启动服务
    stop             停止服务
    status           查看服务状态
  version, -v        显示版本
  help, -h           显示帮助

示例:
  lanlink                      # 启动服务
  lanlink install              # 安装到系统PATH
  lanlink service install      # 安装为系统服务（开机自启）
  lanlink status               # 查看状态
  lanlink list --online        # 仅显示在线节点
  lanlink logs -f              # 实时查看日志
  lanlink ping server.local    # 测试连接

更多信息: https://github.com/618lf/lanlink
`, color(ColorBold, "LanLink - 局域网域名自动映射工具"))
}

