package logger

import (
	"log"
	"os"
)

// InfoLogger 信息级别日志输出到标准输出
var (
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
)

// init 初始化日志器，InfoLogger 输出到 stdout，ErrorLogger 输出到 stderr
func init() {
	InfoLogger = log.New(os.Stdout, "[INFO] ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(os.Stderr, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile)
}

// Info 输出信息级别日志
func Info(format string, v ...interface{}) {
	InfoLogger.Printf(format, v...)
}

// Error 输出错误级别日志
func Error(format string, v ...interface{}) {
	ErrorLogger.Printf(format, v...)
}
