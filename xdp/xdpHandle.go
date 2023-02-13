package xdp

import "os"

func ListenExit() {
	logger.Println("开始监听系统退出信号...")
	for {
		select {
		case <-CtrlC:
			DetachIfaceXdp()
			os.Exit(0)
		}
	}
}
