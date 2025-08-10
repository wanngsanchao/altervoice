package main

import (
	"fmt"
    "log"
    "log/syslog"
	"os"
	"path/filepath"
	"runtime"
)

//define the log level
const (
    LEVEL_WARN = 0
    LEVEL_ERROR = 1
    LEVEL_FATAL = 2
    LEVEL_INFO = 3

)
// define the color of the level
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"    // 红色：错误
	colorYellow = "\033[33m"    // 黄色：警告
//	colorGreen  = "\033[32m"    // 绿色：成功
	colorCyan   = "\033[36m"    // 青色：信息
	colorBold   = "\033[1m"     // 加粗
)

func InitLog(daemon *dae) error{
    if daemon != nil {
        logfile,err := syslog.New(syslog.LOG_EMERG | syslog.LOG_LOCAL5,"")
        if (err != nil) {
            return err
        }

        log.SetPrefix("[altervoice_process] ")
        log.SetOutput(logfile)
        return nil
    }else {
        log.SetPrefix("[altervoice_process] ")
        log.SetOutput(os.Stdout)
        
        return nil
    }
}

func LogPrint(level int,format string,v...interface{}) {
    _,file,line,_ := runtime.Caller(1)
    file = filepath.Base(file)

    switch level {
        case LEVEL_WARN:
            log.Printf(fmt.Sprintf("%s%s[WARN]%s (%s:%d) %s",colorYellow,colorBold,colorReset,file,line,format),v...)
        case LEVEL_ERROR:
            log.Printf(fmt.Sprintf("%s%s[ERROR]%s (%s:%d) %s",colorRed,colorBold,colorReset,file,line,format),v...)
        case LEVEL_FATAL:
            log.Printf(fmt.Sprintf("%s[FATAL]%s (%s:%d) %s",colorRed,colorBold,colorReset,file,line,format),v...)
        case LEVEL_INFO:
            log.Printf(fmt.Sprintf("%s[INFO]%s (%s:%d) %s",colorCyan,colorBold,colorReset,file,line,format),v...)
        default:
            log.Printf(format,v...)

    }
} 
