package services

import (
	"sync"
	"utils"
	"rpcs"
	"dao"
)

type heightService struct {
	BaseService
	sync.Once
}

var _heightService *heightService

func GetHeightService() *heightService {
	if _heightService == nil {
		_heightService = new(heightService)
		_heightService.Once = sync.Once {}
		_heightService.Once.Do(func() {
			_heightService.create()
		})
	}
	return _heightService
}

func (service *heightService) BeforeTurn(s *utils.Status, tgtStt int) {
	switch tgtStt {
	case INIT:
		utils.LogMsgEx(utils.INFO, "initialization", nil)
	case START:
		utils.LogMsgEx(utils.INFO, "start", nil)
	}
}

func (service *heightService) AfterTurn(s *utils.Status, srcStt int) {
	switch s.Current() {
	case INIT:
		utils.LogMsgEx(utils.INFO, "initialized", nil)
	case START:
		// 启动子协程同步进度表的高度
		go service.syncProcessHeight()
		utils.LogMsgEx(utils.INFO, "started", nil)
	}
}

func (service *heightService) syncProcessHeight() {
	asset := utils.GetConfig().GetCoinSettings().Name
	rpc := rpcs.GetRPC(asset)
	var err error
	for err == nil && service.status.Current() == START {
		var curHeight uint64
		if curHeight, err = rpc.GetCurrentHeight(); err != nil {
			utils.LogMsgEx(utils.ERROR, "获取块高失败：%v", err)
			continue
		}
		if _, err = dao.GetProcessDAO().UpdateHeight(asset, curHeight); err != nil {
			utils.LogMsgEx(utils.ERROR, "更新块高失败：%v", err)
			continue
		}
	}
	service.status.TurnTo(DESTORY)
}