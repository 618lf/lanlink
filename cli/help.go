package cli

import (
	"fmt"
)

// ShowHelp 显示帮助信息
func ShowHelp() {
	fmt.Printf(`
%s

用法:
  lanlink [选项]

选项:
  (无参数)          启动服务（前台运行）
  -d, --daemon      安装为系统服务并启动（开机自启）
  -s, --status      查看运行状态
      --stop        停止服务
      --uninstall   卸载系统服务
  -v, --version     显示版本信息
  -h, --help        显示帮助信息

示例:
  lanlink                # 启动服务（前台运行）
  lanlink --daemon       # 安装为后台服务（开机自启）
  lanlink --status       # 查看状态
  lanlink --stop         # 停止服务
  lanlink --uninstall    # 卸载服务

说明:
  LanLink 启动后会自动：
  - 基于硬件序列号生成唯一域名（如 win-abc123.coobee.local）
  - 通过组播发现局域网内其他 LanLink 节点
  - 自动更新 hosts 文件，实现域名到 IP 的映射

  首次运行需要管理员权限（修改 hosts 文件）

更多信息: https://github.com/618lf/lanlink
`, color(ColorBold, "LanLink - 局域网域名自动映射工具"))
}
