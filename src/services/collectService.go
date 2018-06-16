package services

import (
	"sync"
	"utils"
)

type collectService struct {
	sync.Once
	status utils.Status
}

var _collectService *collectService

func GetCollectService() *collectService {
	if _collectService == nil {
		_collectService = new(collectService)
		_collectService.Once = sync.Once {}
		_collectService.Once.Do(func() {
			_collectService.create()
		})
	}
	return _collectService
}

func (service *collectService) create() error {
	service.status.RegAsObs(service)
	service.status.Init([]int { DESTORY, CREATE, INIT, START, STOP })
	return nil
}

func (service *collectService) BeforeTurn(s *utils.Status, tgtStt int) {
	var err error
	switch tgtStt {
	case INIT:
		utils.LogMsgEx(utils.INFO, "initialization", nil)
	case START:
		utils.LogMsgEx(utils.INFO, "start", nil)
	}
}

func (service *collectService) AfterTurn(s *utils.Status, srcStt int) {
	switch s.Current() {
	case INIT:
		utils.LogMsgEx(utils.INFO, "initialized", nil)
	case START:
		utils.LogMsgEx(utils.INFO, "started", nil)
	}
}