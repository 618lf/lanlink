package cli

import (
	"fmt"
	"runtime"
)

const (
	Version   = "1.0.0"
	BuildDate = "2024-11-27"
)

// ShowVersion 显示版本信息
func ShowVersion() {
	fmt.Printf(`
LanLink v%s
Build: %s
Go Version: %s
Platform: %s/%s

局域网域名自动映射工具
项目地址: https://github.com/618lf/lanlink
`, Version, BuildDate, runtime.Version(), runtime.GOOS, runtime.GOARCH)
}

