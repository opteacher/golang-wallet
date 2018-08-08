package main

import (
	"utils"
	"services"
	"unsafe"
	"apis"
	"log"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"net/http"
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

func stopServices() {
	for _, svc := range services.GetInitedServices() {
		svc.Stop()
	}
}

func safeExit() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	<- c
	utils.LogMsgEx(utils.WARNING, "正在安全退出", nil)

	svcStts := make(map[string]int)
	for notYet := true; notYet; {
		for _, svc := range services.GetInitedServices() {
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

	// 关闭日志的存储介质
	utils.CloseAllLogStorage()
	log.Println("退出完毕")
}

var depositServices = []*services.BaseService {
	(*services.BaseService)(unsafe.Pointer(services.GetDepositService())),
	(*services.BaseService)(unsafe.Pointer(services.GetStableService())),
}

var collectServices = []*services.BaseService {
	(*services.BaseService)(unsafe.Pointer(services.GetCollectService())),
}

var withdrawServices = []*services.BaseService {
	(*services.BaseService)(unsafe.Pointer(services.GetHeightService())),
	(*services.BaseService)(unsafe.Pointer(services.GetWithdrawService())),
	(*services.BaseService)(unsafe.Pointer(services.GetStableService())),
}

func main() {
	bsSet := utils.GetConfig().GetBaseSettings()
	sbSet := utils.GetConfig().GetSubsSettings()

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
	if sbSet.APIs.RPC.Active {
		go func() {
			port := 0
			strPort := os.Getenv("PORT")
			if len(strPort) != 0 {
				port, _ = strconv.Atoi(os.Getenv("PORT"))
			} else {
				port = sbSet.APIs.Socket.Port
			}
			utils.LogMsgEx(utils.INFO, "HTTP服务器监听于：%d", port)
			http.HandleFunc("/", apis.HttpHandler)
			log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))

			utils.LogMsgEx(utils.WARNING, "正在安全退出", nil)
			stopServices()
		}()
	}
	if sbSet.APIs.Socket.Active {
		go func() {
			utils.LogMsgEx(utils.INFO, "SOCKET服务器监听于：%d", sbSet.APIs.Socket.Port)
			url := fmt.Sprintf("localhost:%d", sbSet.APIs.Socket.Port)
			var socket net.Listener
			var err error
			if socket, err = net.Listen("tcp", url); err != nil {
				log.Fatal(err)
			}
			defer socket.Close()
			for {
				var conn net.Conn
				if conn, err = socket.Accept(); err != nil {
					utils.LogMsgEx(utils.WARNING, "接收SOCKET信息错误：%v", err)
					continue
				}
				utils.LogMsgEx(utils.INFO, "接收来自：%s的消息", conn.RemoteAddr().String())
				apis.SocketHandler(conn)
			}
			utils.LogMsgEx(utils.WARNING, "正在安全退出", nil)
			stopServices()
		}()
	}

	// 处理服务安全退出
	safeExit()
}