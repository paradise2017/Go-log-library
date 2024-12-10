/*
1: 第三方包不需要main
2：结构体包含构造函数
3：结构体方法
*/
package mylogger

import (
	"fmt"
	"os"
	"path"
	"time"
)

// FileLogger文件记录日志
type FileLogger struct {
	level         Level // 日志级别，供测试和开发使用
	log_file_path string
	log_file_name string
	log_file      *os.File
	err_file      *os.File
	max_size      int64
}

// FileLogger:构造函数  首字母大写对外开发
func NewFileLogger(level Level, log_file_path, log_file_name string) *FileLogger {
	file_obj := &FileLogger{
		level,
		log_file_path,
		log_file_name,
		nil,
		nil,
		1 * 1024,
	}
	// 初始化日志句柄
	file_obj.initFileLogger()
	return file_obj
}

// 初始化文件句柄
func (f *FileLogger) initFileLogger() {
	// 1:打开日志信息文件
	filepath := path.Join(f.log_file_path, f.log_file_name)
	//filepath := fmt.Sprintf("%s/%s", f.log_file_path, f.log_file_name)
	file_obj, err := os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0755)
	if err != nil {
		panic(fmt.Sprintf("open file failed"))
	}
	f.log_file = file_obj

	// 2:打开错误信息文件
	err_log_name := fmt.Sprintf("%s.err", filepath)
	err_file_obj, err := os.OpenFile(err_log_name, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0755)
	if err != nil {
		panic(fmt.Sprintf("open file failed"))
	}
	f.err_file = err_file_obj
}

// 检查是否拆分
func (f *FileLogger) checkSplit(file_obj *os.File) bool {
	file_info, _ := file_obj.Stat()
	file_size := file_info.Size()
	return file_size >= f.max_size
}

// 封装切分日志
func (f *FileLogger) splitLogFile(file_obj **os.File) {

	// 切分文件 1：origin关闭，2：备份原来的 3：新建一个文件
	file_name := (*file_obj).Name()
	back_up_name := fmt.Sprintf("%s_%v.back", file_name, time.Now().Unix())
	(*file_obj).Close()

	os.Rename(file_name, back_up_name) // 原来的更新名字

	file_obj_temp, err := os.OpenFile(file_name, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0755) // 拿到新的名字注册
	if err != nil {
		panic(fmt.Errorf("open log failed"))
	}
	*file_obj = file_obj_temp
}

// 记录日志核心
func (f *FileLogger) log(level Level, format string, args ...any) {
	if f.level > level { // 日志级别大，不用写日志
		return
	}
	// 标准化 [2006-01-02 15:04:05.000]
	now_str := time.Now().Format("[2006-01-02 15:04:05.000]")
	file_name, line, func_name := getCallerInfo()
	format = fmt.Sprintf("%s [%s] [file:%s;func:%s] [line:%d] %s", now_str, getLevelStr(f.level), file_name, func_name, line, format)

	// 写日志之前检查，检查当前日志文件大小
	if f.checkSplit(f.log_file) {
		f.splitLogFile(&f.log_file)
	}
	fmt.Fprintf(f.log_file, format, args...) // 将日志信息导入 log_file句柄

	fmt.Fprintln(f.log_file)

	if level >= ERRORLevel {
		if f.checkSplit(f.err_file) {
			f.splitLogFile(&f.err_file)
		}

		fmt.Fprintf(f.err_file, format, args...)
	}
}

// 区别printf ，Fprintf输出到任意文件
// 记录日志 任意类 format string, args ...any
func (f *FileLogger) Debug(format string, args ...any) {
	f.log(DEBUGLevel, format, args...)
}

// TRACELevel

func (f *FileLogger) Trace(format string, args ...any) {
	f.log(TRACELevel, format, args...)
}

// Info
func (f *FileLogger) Info(format string, args ...any) {
	f.log(INFOLevel, format, args...)
}

// Warning
func (f *FileLogger) Warn(format string, args ...any) {
	f.log(WARNLevel, format, args...)
}

// Error
func (f *FileLogger) Error(format string, args ...any) {
	f.log(ERRORLevel, format, args...)
}

// CIRTALLevel

func (f *FileLogger) Critical(format string, args ...any) {
	f.log(CIRTALLevel, format, args...)
}

// Close 关闭日志文件句柄
func (f *FileLogger) Close() {
	f.log_file.Close()
	f.err_file.Close()
}
