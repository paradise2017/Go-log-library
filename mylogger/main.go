package mylogger

// 自定义日志库，实现日志记录功能

// 定义具体的日志级别
type Level uint16

const (
	DEBUGLevel Level = iota // 不写默认相同，下面常量自增
	TRACELevel
	INFOLevel
	WARNLevel
	ERRORLevel
	CIRTALLevel
)

// logger方法的抽象接口
type Logger interface {
	DEBUG(format string, args ...any)
	Info(format string, args ...any)
	Trace(format string, args ...any)
	Warn(format string, args ...any)
	Error(format string, args ...any)
	Critical(format string, args ...any)
	Close()
}

func getLevelStr(level Level) string {
	switch level {
	case 0:
		return "DEBUG"
	case 1:
		return "TRACE"
	case 2:
		return "INFO"
	case 3:
		return "WARN"
	case 4:
		return "ERROR"
	case 5:
		return "CIRTAL"
	default:
		return "DEBUG"
	}
}
