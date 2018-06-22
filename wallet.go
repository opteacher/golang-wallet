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

const (
	NONE = iota
	HELP
	SERVICE
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

func main() {
	var runStt utils.Status
	runStt.Init([]int {
		NONE, HELP, SERVICE,
	})
	argNum := len(os.Args) - 1
	cmdSet := utils.GetConfig().GetCmdsSettings()
	subSet := utils.GetConfig().GetSubsSettings()
	for i, arg := range os.Args[1:] {
		arg = strings.ToLower(arg)
		switch arg {
		case "help":
			runStt.TurnTo(HELP)
			if i == argNum - 1 {
				fmt.Println(cmdSet.Help)
			}
		case "version":
			fmt.Println(cmdSet.Version)
		case "test":
			fmt.Println("测试")
		case "service":
			runStt.TurnTo(SERVICE)
		default:
			switch runStt.Current() {
			case HELP:
				fmt.Println("帮助：" + arg)
			case SERVICE:
				switch arg {
				case "deposit":
					runDeposit()
				case "collect":
					runCollect()
				case "withdraw":
					http.HandleFunc("/api/withdraw", apis.WithdrawHandler)
				default:
					fmt.Println("指定的服务不存在")
				}
				utils.LogMsgEx(utils.INFO, "服务器监听于：%d", subSet.Server.Port)
				log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", subSet.Server.Port), nil))
			default:
				fmt.Println(cmdSet.Unknown)
			}
			runStt.TurnTo(NONE)
		}
	}
}