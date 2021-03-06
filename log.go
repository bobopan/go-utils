package baseutils

import (
	"container/list"
	"os"
	"time"

	"github.com/op/go-logging"
)

//LogInterface log日志接口，实现get path和get log level方法
type LogInterface interface {
	GetLogPath() string
	GetLogLevel() string
	GetLogErrLevel() string
}

type defaultLogger struct {
}

func (l defaultLogger) GetLogPath() string {
	return ""
}

func (l defaultLogger) GetLogLevel() string {
	return "DEBUG"
}

func (l defaultLogger) GetLogErrLevel() string {
	return "WARNING"
}

var (
	//Log 日志
	Log          = logging.MustGetLogger("mgtv")
	fileList     = list.New()
	logInterface LogInterface
)

// Example format string. Everything except the message has a custom color
// which is dependent on the log level. Many fields have a custom output
// formatting too, eg. the time returns the hour down to the milli second.
//     %{id}        Sequence number for log message (uint64).
//     %{pid}       Process id (int)
//     %{time}      Time when log occurred (time.Time)
//     %{level}     Log level (Level)
//     %{module}    Module (string)
//     %{program}   Basename of os.Args[0] (string)
//     %{message}   Message (string)
//     %{longfile}  Full file name and line number: /a/b/c/d.go:23
//     %{shortfile} Final file name element and line number: d.go:23
//     %{color}     ANSI color based on log level
//     %{longpkg}   Full package path, eg. github.com/go-logging
//     %{shortpkg}  Base package path, eg. go-logging
//     %{longfunc}  Full function name, eg. littleEndian.PutUint32
//     %{shortfunc} Base function name, eg. PutUint32
var stdFormat = logging.MustStringFormatter(
	"%{color}%{time:15:04:05.000} %{shortfile} >%{level:.5s}%{color:reset} - %{message}",
)

var fileFormat = logging.MustStringFormatter(
	"%{time:15:04:05.000} %{shortfile} >%{level:.5s} - %{message}",
)

//关闭旧log 打开的文件
//newFile 本次是否打开了新文件
func closeOldLogFd(newFile bool) {
	expectedFdNum := 0
	if newFile {
		expectedFdNum++
	}
	Log.Notice("in closeOld LogFd expected fd:%d, list len:%d", expectedFdNum, fileList.Len())
	if fileList.Len() > expectedFdNum {
		element := fileList.Front()
		if element == nil {
			return
		}
		if fp, ok := element.Value.(*os.File); ok {
			fileList.Remove(element)
			time.Sleep(time.Second * 5)
			Log.Notice("start close old log file")
			fp.Close()
		} else {
			Log.Error("fd type error")
		}
	}
}

// 如果 path 路径不为空，则使用文件记录日志
// 否则使用std out 输出日志
// change by zzh 20151130
// SetBackend  可重复调用
func initLog(path string, level logging.Level, errLevel logging.Level) error {
	if len(path) > 0 {
		fp, err := os.OpenFile(path, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
		if err != nil {
			return err
		}
		fileList.PushBack(fp)
		fileBackend := logging.NewLogBackend(fp, "", 1)
		fileFormatter := logging.NewBackendFormatter(fileBackend, fileFormat)
		fileB := logging.AddModuleLevel(fileFormatter)
		fileB.SetLevel(level, "")

		//warning级别的日志需要在控制台的，ERR输出
		backendErr := logging.NewLogBackend(os.Stderr, "", 1)
		backendErrFormatter := logging.NewBackendFormatter(backendErr, stdFormat)
		stdErr := logging.AddModuleLevel(backendErrFormatter)
		stdErr.SetLevel(errLevel, "")
		logging.SetBackend(fileB, stdErr)
	} else {
		stdBackend := logging.NewLogBackend(os.Stdout, "", 1)
		stdFormatter := logging.NewBackendFormatter(stdBackend, stdFormat)
		stdB := logging.AddModuleLevel(stdFormatter)
		stdB.SetLevel(level, "")

		//warning级别的日志需要在控制台的，ERR输出
		backendErr := logging.NewLogBackend(os.Stderr, "", 1)
		backendErrFormatter := logging.NewBackendFormatter(backendErr, stdFormat)
		stdErr := logging.AddModuleLevel(backendErrFormatter)
		stdErr.SetLevel(errLevel, "")

		logging.SetBackend(stdB, stdErr)
	}
	go closeOldLogFd(len(path) > 0)
	return nil
}

func reloadLog() error {
	logPath := logInterface.GetLogPath()
	logLevel := logInterface.GetLogLevel()
	level, err := logging.LogLevel(logLevel)
	if err != nil {
		level = logging.INFO
	}
	errLevelStr := logInterface.GetLogErrLevel()
	errLevel, err := logging.LogLevel(errLevelStr)
	if err != nil {
		errLevel = logging.WARNING
	}
	if err := initLog(logPath, level, errLevel); err != nil {
		return err
	}
	return nil
}

// InitLog 日志初始化
func InitLog(log LogInterface) error {
	if log != nil {
		logInterface = log
	}
	return reloadLog()
}

func init() {
	logInterface = defaultLogger{}
}
