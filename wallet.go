package main

import (
	"os"
	"utils"
	"fmt"
	"strings"
	"services"
	"os/signal"
	"unsafe"
	"log"
	"net/http"
	"apis"
)

func runService(svcs []*services.BaseService) {
	var svc *services.BaseService
	for _, svc = range svcs {
		svc.Init()
	}

	for _, svc = range svcs {
		svc.Start()
	}

	c := make(chan os.Signal, 1)
	stop := make(chan bool)
	signal.Notify(c, os.Interrupt, os.Kill)

	go func() {
		<- c
		utils.LogMsgEx(utils.WARNING, "正在安全退出", nil)
		for _, svc = range svcs {
			svc.Stop()
		}
		stop <- true
	}()

	<- stop
	for notYet := true; notYet; {
		for _, svc = range svcs {
			if !svc.IsDestroy() {
				notYet = true
				break
			} else {
				notYet = false
			}
		}
	}
	utils.LogMsgEx(utils.INFO, "退出完毕", nil)
}

func runDeposit() {
	runService([]*services.BaseService {
		(*services.BaseService)(unsafe.Pointer(services.GetDepositService())),
		(*services.BaseService)(unsafe.Pointer(services.GetNotifyService())),
	})
}

func runCollect() {
	runService([]*services.BaseService {
		(*services.BaseService)(unsafe.Pointer(services.GetCollectService())),
	})
}

func runWithdraw() {

}

func main() {
	cmdSet := utils.GetConfig().GetCmdsSettings()
	subSet := utils.GetConfig().GetSubsSettings()
	// 收集参数
	args := make(map[string]string)
	curArg := ""
	for _, arg := range os.Args[1:] {
		if arg[0] == '-' {
			curArg = arg[1:]
			args[curArg] = ""
		} else {
			args[curArg] = arg
			curArg = ""
		}
	}
	// 根据参数启动相应的配置
	for key, val := range args {
		switch strings.ToLower(key) {
		case "help":
			if val == "" {
				fmt.Println(cmdSet.Help)
			} else {
			}
		case "version":
			fmt.Println(cmdSet.Version)
		case "service":
			switch strings.ToLower(val) {
			case "deposit":
				runDeposit()
			case "collect":
				runCollect()
			case "withdraw":

			default:

			}
		case "rpc":
			switch strings.ToLower(val) {
			case "http":
				fallthrough
			default:
				utils.LogMsgEx(utils.INFO, "服务器监听于：%d", subSet.Server.Port)
				http.HandleFunc("/", apis.RootHandler)
				log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", subSet.Server.Port), nil))
			}
		}
	}
}