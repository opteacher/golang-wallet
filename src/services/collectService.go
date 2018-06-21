package services

import (
	"sync"
	"utils"
	"rpcs"
	"dao"
	"time"
)

type collectService struct {
	sync.Once
	status utils.Status
	addresses []string
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
		go service.doCollect()
		utils.LogMsgEx(utils.INFO, "started", nil)
	}
}

func (service *collectService) doCollect() {
	var err error
	coinSet := utils.GetConfig().GetCoinSettings()
	rpc := rpcs.GetRPC(coinSet.Name)
	for err == nil && service.status.Current() == START {
		addrAmount := make(map[string]float64)
		if addrAmount, err = rpc.GetDepositAmount(); err != nil {
			utils.LogMsgEx(utils.ERROR, "获取充值地址上余额时，发生错误：%v", err)
			continue
		}

		for addr, balance := range addrAmount {
			var txHash string
			if txHash, err = rpc.SendFrom(addr, coinSet.Collect, balance); err != nil {
				utils.LogMsgEx(utils.ERROR, "发送归集请求失败：%v", err)
				err = nil
				continue
			}

			if _, err = dao.GetCollectDAO().AddSentCollect(txHash, coinSet.Name, addr, balance); err != nil {
				utils.LogMsgEx(utils.ERROR, "添加发送的归集记录失败：%v", err)
				continue
			}
		}

		time.Sleep(coinSet.CollectInterval * time.Second)
	}
}