package baseutils

import (
	//"config"
	"os"
	"os/signal"
	"syscall"
)

//SignalReload interface
type SignalReload interface {
	Reload()
}

//InitSignal 用户信号量初始化
//
// SIGUSR1: 日志文件重新打开类似 nginx -s reload, 完成日志切割
//
// SIGUSR2: 配置文件重新加载,可以完成比如日志级别的动态改变
func InitSignal(signalReload SignalReload) {
	s1 := make(chan os.Signal, 1)
	s2 := make(chan os.Signal, 1)
	signal.Notify(s1, syscall.SIGUSR1)
	signal.Notify(s2, syscall.SIGUSR2)
	go func() {
		for {
			select {
			case _ = <-s1:
				NewContext("redis").Notice("received signal USR1")
				reloadLog()
			case _ = <-s2:
				NewContext("redis").Notice("received signal USR2")
				if signalReload != nil {
					signalReload.Reload()
				}
				reloadLog()
			}
		}
	}()
}
