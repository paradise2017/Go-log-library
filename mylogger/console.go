package mylogger

import (
	"fmt"
	"os"
	"time"
)

// 终端打印日志

// ConsoleLogger终端日志
type ConsoleLogger struct {
	level    Level
	file     *os.File
	err_file *os.File
}

// ConsoleLogger:构造函数  首字母大写对外开发
func NewConsoleLogger(console_level Level) *ConsoleLogger {
	return &ConsoleLogger{
		console_level,
		os.Stdout,
		os.Stderr,
	}
}

// 记录日志核心
func (c *ConsoleLogger) log(level Level, format string, args ...any) {
	if c.level > level { // 日志级别大，不用写日志
		return
	}
	// 标准化 [2006-01-02 15:04:05.000]
	now_str := time.Now().Format("[2006-01-02 15:04:05.000]")
	file_name, line, func_name := getCallerInfo()
	format = fmt.Sprintf("%s [%s] [file:%s;func:%s] [line:%d] %s", now_str, getLevelStr(c.level), file_name, func_name, line, format)
	fmt.Fprintf(c.file, format, args...) // 将日志信息导入 Console句柄
	fmt.Fprintln(c.file)

	// 错误信息
	if c.level >= ERRORLevel {
		fmt.Fprintf(c.err_file, format, args...)
		fmt.Fprintln(c.err_file)
	}
}

// Debug
func (c *ConsoleLogger) ConsoleDebug(format string, args ...any) {
	c.log(DEBUGLevel, format, args...)
}

// TRACELevel

func (c *ConsoleLogger) ConsoleTrace(format string, args ...any) {
	c.log(TRACELevel, format, args...)
}

// Info
func (c *ConsoleLogger) ConsoleInfo(format string, args ...any) {
	c.log(INFOLevel, format, args...)
}

// Warning
func (c *ConsoleLogger) ConsoleWarn(format string, args ...any) {
	c.log(WARNLevel, format, args...)
}

// Error
func (f *FileLogger) ConsoleError(format string, args ...any) {
	f.log(ERRORLevel, format, args...)
}

// ConsoleCritical

func (c *ConsoleLogger) ConsoleCritical(format string, args ...any) {
	c.log(CIRTALLevel, format, args...)
}

// Close 关闭日志文件句柄
func (c *ConsoleLogger) Close() {

}
