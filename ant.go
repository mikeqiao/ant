package ant

import (
	"os"
	"os/signal"

	"github.com/mikeqiao/ant/group"
	"github.com/mikeqiao/ant/log"
	mod "github.com/mikeqiao/ant/module"
	"github.com/mikeqiao/ant/network"
)

func Start(f Func) {
	//初始化基本设置
	Init()
	//初始化功能设置
	f.Init()
	//开始运行服务程序
	Run()
	//设置信号接收
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	sig := <-c
	log.Release("bee closing down (signal: %v)", sig)
	//关闭服务
	Close()
	//关闭添加功能
	f.Close()
	//等待所以线程结束
	group.Wait()
}

func Init() {
	log.Init()
	mod.Init()
	network.Init()
}

func Run() {
	log.Run()
	mod.Run()
	network.Run()
}

func Close() {
	network.Close()
	mod.Close()
	log.Close()
}
