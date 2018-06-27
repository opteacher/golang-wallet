package main

import (
	"os"
	"utils"
	"fmt"
	"strings"
	"services"
	"unsafe"
	"log"
	"net/http"
	"apis"
	"os/signal"
)

func initServices(svcs []*services.BaseService) {
	for _, svc := range svcs {
		svc.Init()
	}
}

func runServices() {
	for _, svc := range services.GetInitedServices() {
		svc.Start()
	}
}

func exitCtrl() {
	svcs := services.GetInitedServices()

	c := make(chan os.Signal, 1)
	stop := make(chan bool)
	signal.Notify(c, os.Interrupt, os.Kill)

	go func() {
		<- c
		utils.LogMsgEx(utils.WARNING, "正在安全退出", nil)
		for _, svc := range svcs {
			svc.Stop()
		}
		stop <- true
	}()

	<- stop
	svcStts := make(map[string]int)
	for notYet := true; notYet; {
		for _, svc := range svcs {
			if !svc.IsDestroy() {
				if stt, ok := svcStts[svc.Name()]; ok {
					if stt == svc.CurrentStatus() { continue }
				}
				utils.LogMsgEx(utils.WARNING, "%s服务还未安全退出，处于状态：%s",
					svc.Name(), services.ServiceStatus[svc.CurrentStatus()])
				svcStts[svc.Name()] = svc.CurrentStatus()
				notYet = true
				break
			} else {
				notYet = false
			}
		}
	}
	utils.LogMsgEx(utils.INFO, "退出完毕", nil)
}

var depositServices = []*services.BaseService {
	(*services.BaseService)(unsafe.Pointer(services.GetDepositService())),
	(*services.BaseService)(unsafe.Pointer(services.GetNotifyService())),
}

var collectServices = []*services.BaseService {
	(*services.BaseService)(unsafe.Pointer(services.GetCollectService())),
}

var withdrawServices = []*services.BaseService {
	(*services.BaseService)(unsafe.Pointer(services.GetWithdrawService())),
	(*services.BaseService)(unsafe.Pointer(services.GetNotifyService())),
}

func main() {
	cmdSet := utils.GetConfig().GetCmdsSettings()
	subSet := utils.GetConfig().GetSubsSettings()
	// 收集参数
	args := make(map[string]string)
	curArg := ""
	for _, arg := range os.Args[1:] {
		if arg[0] == '-' {
			curArg = strings.ToLower(arg[1:])
			args[curArg] = ""
		} else {
			args[curArg] = arg
			curArg = ""
		}
	}
	// 根据参数启动相应的配置
	for key, val := range args {
		switch key {
		case "help":
			if val == "" {
				fmt.Println(cmdSet.Help)
			} else {
			}
			return
		case "version":
			fmt.Println(cmdSet.Version)
			return
		case "service":
			switch strings.ToLower(val) {
			case "deposit":
				initServices(depositServices)
			case "collect":
				initServices(collectServices)
			case "withdraw":
				initServices(withdrawServices)
			case "all":
				initServices(depositServices)
				initServices(collectServices)
				initServices(withdrawServices)
			default:
				utils.LogMsgEx(utils.WARNING, "找不到指定的服务：%s", val)
			}
		case "remote":
			switch strings.ToLower(val) {
			case "http":
				fallthrough
			default:
				go func() {
					utils.LogMsgEx(utils.INFO, "服务器监听于：%d", subSet.Server.Port)
					http.HandleFunc("/", apis.RootHandler)
					log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", subSet.Server.Port), nil))
				}()
			}
		}
	}
	// 启动服务器
	runServices()
	exitCtrl()
}