// util包用于提供各种实用工具函数
package util

// 导入所需的包
import (
	// 导入操作系统功能包
	"os"
	// 导入信号处理包
	"os/signal"
	// 导入系统调用包
	"syscall"
)

// WaitQuit 等待程序退出信号，并执行清理函数
func WaitQuit(fns ...func()) {
	// 创建一个信号通道
	c := make(chan os.Signal, 1)
	// 监听中断和终止信号
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	// 阻塞等待信号
	<-c
	// 依次执行所有传入的清理函数
	for _, fn := range fns {
		fn()
	}
}
