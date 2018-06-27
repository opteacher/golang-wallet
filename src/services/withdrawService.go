package services

import (
	"sync"
	"utils"
)

/***
	提币服务：接收来自API的提币操作，并启动子协程等待其入链后交给通知服务
	子协程（startWaitForInchain）：查询发起的提币交易，检查是否入链
 */
type withdrawService struct {
	BaseService
	sync.Once
	addresses []string
	height uint64
}

var _withdrawService *withdrawService

func GetWithdrawService() *withdrawService {
	if _withdrawService == nil {
		_withdrawService = new(withdrawService)
		_withdrawService.Once = sync.Once {}
		_withdrawService.Once.Do(func() {
			_withdrawService.create()
		})
	}
	return _withdrawService
}

func (service *withdrawService) create() error {
	service.status.RegAsObs(service)
	return service.BaseService.create()
}

func (service *withdrawService) BeforeTurn(s *utils.Status, tgtStt int) {
	switch tgtStt {
	case INIT:
		utils.LogMsgEx(utils.INFO, "initialization", nil)
	case START:
		utils.LogMsgEx(utils.INFO, "start", nil)
	}
}

func (service *withdrawService) AfterTurn(s *utils.Status, srcStt int) {
	switch s.Current() {
	case INIT:
		utils.LogMsgEx(utils.INFO, "initialized", nil)
	case START:
		utils.LogMsgEx(utils.INFO, "started", nil)
	}
}