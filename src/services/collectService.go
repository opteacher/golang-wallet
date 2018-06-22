package services

import (
	"sync"
	"utils"
	"rpcs"
	"dao"
	"time"
)

type collectService struct {
	BaseService
	sync.Once
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
	return service.BaseService.create()
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
		// 启动归集协程
		go service.doCollect()
		utils.LogMsgEx(utils.INFO, "started", nil)
	}
}

func (service *collectService) doCollect() {
	var err error
	coinSet := utils.GetConfig().GetCoinSettings()
	rpc := rpcs.GetRPC(coinSet.Name)
	for err == nil && service.status.Current() == START {
		utils.LogMsgEx(utils.INFO, "开始一轮归集", nil)

		// 获取所有有余额需要归集的地址/账户
		addrAmount := make(map[string]float64)
		if addrAmount, err = rpc.GetDepositAmount(); err != nil {
			utils.LogMsgEx(utils.ERROR, "获取充值地址上余额时，发生错误：%v", err)
			continue
		}
		utils.LogMsgEx(utils.INFO, "发现%d个地址/账户需要归集", len(addrAmount))

		for addr, balance := range addrAmount {
			// 发起转账请求
			var txHash string
			if txHash, err = rpc.SendFrom(addr, coinSet.Collect, balance); err != nil {
				utils.LogMsgEx(utils.ERROR, "发送归集请求失败：%v", err)
				err = nil
				continue
			}
			utils.LogMsgEx(utils.INFO, "地址/账户%s下的余额已被归集", addr)
			utils.LogMsgEx(utils.INFO, "归集的交易ID为：%s", txHash)

			// 把归集的记录保存进数据库
			if _, err = dao.GetCollectDAO().AddSentCollect(txHash, coinSet.Name, addr, balance); err != nil {
				utils.LogMsgEx(utils.ERROR, "添加发送的归集记录失败：%v", err)
				continue
			}
			utils.LogMsgEx(utils.INFO, "%s归集记录已持久化到数据库", txHash)
		}

		// 定时
		time.Sleep(coinSet.CollectInterval * time.Second)
	}

	service.status.TurnTo(DESTORY)
}