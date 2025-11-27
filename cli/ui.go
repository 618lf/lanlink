package cli

import (
	"fmt"
	"runtime"
	"syscall"
	"unsafe"
)

// ANSI 颜色代码
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorCyan   = "\033[36m"
	ColorGray   = "\033[90m"
	ColorBold   = "\033[1m"
)

var colorEnabled = true

// InitColors 初始化颜色支持
func InitColors() {
	if runtime.GOOS == "windows" {
		// 启用 Windows 10+ 的 ANSI 颜色支持
		kernel32 := syscall.NewLazyDLL("kernel32.dll")
		setConsoleMode := kernel32.NewProc("SetConsoleMode")
		var mode uint32
		handle := syscall.Handle(syscall.Stdout)
		syscall.Syscall(kernel32.NewProc("GetConsoleMode").Addr(), 2, uintptr(handle), uintptr(unsafe.Pointer(&mode)), 0)
		mode |= 0x0004 // ENABLE_VIRTUAL_TERMINAL_PROCESSING
		setConsoleMode.Call(uintptr(handle), uintptr(mode))
	}
}

func init() {
	InitColors()
}

// 辅助函数
func color(c string, s string) string {
	if !colorEnabled {
		return s
	}
	return c + s + ColorReset
}

// Success 成功消息
func Success(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Printf("%s✓%s %s\n", ColorGreen, ColorReset, msg)
}

// Error 错误消息
func Error(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Printf("%s✗%s %s\n", ColorRed, ColorReset, msg)
}

// Warn 警告消息
func Warn(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Printf("%s⚠%s %s\n", ColorYellow, ColorReset, msg)
}

// Info 信息消息
func Info(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Println(msg)
}

// Header 标题
func Header(title string) {
	line := "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
	fmt.Printf("\n%s\n", color(ColorCyan, line))
	fmt.Printf("%s  %s\n", color(ColorCyan, ""), color(ColorBold, title))
	fmt.Printf("%s\n\n", color(ColorCyan, line))
}

// Footer 页脚
func Footer() {
	line := "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
	fmt.Printf("%s\n", color(ColorCyan, line))
}

// Section 章节标题
func Section(title string) {
	fmt.Printf("\n%s:\n", color(ColorYellow, title))
}

// KeyValue 键值对输出
func KeyValue(key, value string) {
	fmt.Printf("  %s: %s\n", color(ColorGray, key), value)
}

// StatusIndicator 状态指示器
func StatusIndicator(isOK bool) string {
	if isOK {
		return color(ColorGreen, "✓")
	}
	return color(ColorRed, "✗")
}

