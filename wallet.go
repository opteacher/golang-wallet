package main

import (
	"os"
	"utils"
	"services"
	"unsafe"
	"os/signal"
	"net/http"
	"apis"
	"log"
	"fmt"
)

func initServices(svcs []*services.BaseService) {
	for _, svc := range svcs {
		svc.Init()
		svc.Start()
	}
}

func runServices() {
	for _, svc := range services.GetInitedServices() {
		svc.Start()
	}
}

var depositServices = []*services.BaseService {
	(*services.BaseService)(unsafe.Pointer(services.GetDepositService())),
	(*services.BaseService)(unsafe.Pointer(services.GetNotifyService())),
}

var collectServices = []*services.BaseService {
	(*services.BaseService)(unsafe.Pointer(services.GetCollectService())),
}

var withdrawServices = []*services.BaseService {
	(*services.BaseService)(unsafe.Pointer(services.GetHeightService())),
	(*services.BaseService)(unsafe.Pointer(services.GetWithdrawService())),
	(*services.BaseService)(unsafe.Pointer(services.GetNotifyService())),
}

func main() {
	bsSet := utils.GetConfig().GetBaseSettings()

	// 服务的初始化和启动
	if len(bsSet.Services) == 0 {
		initServices(depositServices)
		initServices(collectServices)
		initServices(withdrawServices)
	} else {
		for _, svc := range bsSet.Services {
			switch svc {
			case "deposit": initServices(depositServices)
			case "collect": initServices(collectServices)
			case "withdraw": initServices(withdrawServices)
			}
		}
	}
	runServices()

	// API的初始化和监听
	svcs := services.GetInitedServices()
	if bsSet.APIs.RPC.Active {
		utils.LogMsgEx(utils.INFO, "服务器监听于：%d", bsSet.APIs.RPC.Port)
		http.HandleFunc("/", apis.RootHandler)
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", bsSet.APIs.RPC.Port), nil))

		utils.LogMsgEx(utils.WARNING, "正在安全退出", nil)
		for _, svc := range svcs { svc.Stop() }
	} else {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, os.Kill)

		go func() {
			<- c
			utils.LogMsgEx(utils.WARNING, "正在安全退出", nil)
			for _, svc := range svcs { svc.Stop() }
		}()
	}

	// 处理服务安全退出
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