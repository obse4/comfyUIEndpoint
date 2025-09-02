package logger

import (
	"log"
	"os"
)

// 定义日志等级类型
type LogLeveL uint32

const (
	DebugLevel LogLeveL = iota //iota=0
	InfoLevel                  //InfoLevel=iota, iota=1
	WarnLevel                  //WarnLevel=iota, iota=2
	ErrorLevel                 //ErrorLevel=iota, iota=3
	FatalLevel
)

var (
	logOut      = os.Stdout
	debugLogger = log.New(logOut, "[DEBUG] ", log.LstdFlags)
	infoLogger  = log.New(logOut, "[INFO] ", log.LstdFlags)
	warnLogger  = log.New(logOut, "[WARN] ", log.LstdFlags)
	errorLogger = log.New(logOut, "[ERROR] ", log.LstdFlags)
	fatalLogger = log.New(logOut, "[FATAL] ", log.LstdFlags)

	//默认的LogLevel为0，即所有级别的日志都打印
	logLevel LogLeveL = 0
)

// 日志配置
type LogConfig struct {
	// 默认的LogLevel为0，即所有级别的日志都打印
	LogLevel LogLeveL `yaml:"logLevel"`
}

// 初始化日志
func InitLogger(conf *LogConfig) {
	conf.setLogLevel()
}

func (l *LogConfig) setLogLevel() {
	logLevel = l.LogLevel
}

func Debug(format string, v ...interface{}) {
	if logLevel <= DebugLevel {
		debugLogger.Printf(format, v...)
	}
}

func Info(format string, v ...interface{}) {
	if logLevel <= InfoLevel {
		infoLogger.Printf(format, v...) //format末尾如果没有换行符会自动加上
	}
}

func Warn(format string, v ...interface{}) {
	if logLevel <= WarnLevel {
		warnLogger.Printf(format, v...)
	}
}

func Error(format string, v ...interface{}) {
	if logLevel <= ErrorLevel {
		errorLogger.Printf(format, v...)
	}
}

func Fatal(format string, v ...interface{}) {
	if logLevel <= FatalLevel {
		fatalLogger.Fatalf(format, v...)
	}
}

func init() {
	InitLogger(&LogConfig{})
}
