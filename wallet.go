package main

import (
	"os"
	"utils"
	"fmt"
	"strings"
	"services"
	"os/signal"
)

const (
	NONE = iota
	HELP
	SERVICE
)

func runDeposit() {
	depositService := services.GetDepositService()
	notifyService := services.GetNotifyService()

	depositService.Init()
	notifyService.Init()

	depositService.Start()
	notifyService.Start()

	c := make(chan os.Signal, 1)
	stop := make(chan bool)
	signal.Notify(c, os.Interrupt, os.Kill)

	go func() {
		<- c
		depositService.Stop()
		notifyService.Stop()
		stop <- true
	}()

	<- stop
	for !depositService.IsDestroy() || !notifyService.IsDestroy() {}
}

func main() {
	var runStt utils.Status
	runStt.Init([]int {
		NONE, HELP, SERVICE,
	})
	argNum := len(os.Args) - 1
	cmdSet := utils.GetConfig().GetCmdsSettings()
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
				case "withdraw":
				default:
					fmt.Println("指定的服务不存在")
				}
			default:
				fmt.Println(cmdSet.Unknown)
			}
			runStt.TurnTo(NONE)
		}
	}
}