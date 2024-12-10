package mylogger

import (
	"path"
	"runtime"
)

func getCallerInfo() (file_path string, line int, func_name string) {
	// 0 当前函数，1代表上一层 ,2 代表上上层函数调用该函数
	pc, file_path, line, ok := runtime.Caller(3)
	if !ok {
		return
	}
	file_path = path.Base(file_path)         // 全路径拿到文件名
	func_name = runtime.FuncForPC(pc).Name() // 根据pc拿到函数名字
	func_name = path.Base(func_name)         // 函数名
	return                                   // 匿名函数返回
}
