package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

// Level 日志级别
type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
)

// Logger 日志记录器
type Logger struct {
	level  Level
	logger *log.Logger
	file   *os.File
}

var std *Logger

// Init 初始化日志
func Init(levelStr string, logFile string) error {
	level := parseLevel(levelStr)

	// 创建日志文件
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	// 同时输出到控制台和文件
	multiWriter := io.MultiWriter(os.Stdout, file)

	std = &Logger{
		level:  level,
		logger: log.New(multiWriter, "", 0),
		file:   file,
	}

	return nil
}

// Close 关闭日志文件
func Close() {
	if std != nil && std.file != nil {
		std.file.Close()
	}
}

// Debug 调试日志
func Debug(format string, v ...interface{}) {
	if std != nil && std.level <= LevelDebug {
		std.log("DEBUG", format, v...)
	}
}

// Info 信息日志
func Info(format string, v ...interface{}) {
	if std != nil && std.level <= LevelInfo {
		std.log("INFO", format, v...)
	}
}

// Warn 警告日志
func Warn(format string, v ...interface{}) {
	if std != nil && std.level <= LevelWarn {
		std.log("WARN", format, v...)
	}
}

// Error 错误日志
func Error(format string, v ...interface{}) {
	if std != nil && std.level <= LevelError {
		std.log("ERROR", format, v...)
	}
}

// log 内部日志方法
func (l *Logger) log(level string, format string, v ...interface{}) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	message := fmt.Sprintf(format, v...)
	l.logger.Printf("[%s] [%s] %s", timestamp, level, message)
}

// parseLevel 解析日志级别
func parseLevel(levelStr string) Level {
	switch strings.ToLower(levelStr) {
	case "debug":
		return LevelDebug
	case "info":
		return LevelInfo
	case "warn":
		return LevelWarn
	case "error":
		return LevelError
	default:
		return LevelInfo
	}
}

