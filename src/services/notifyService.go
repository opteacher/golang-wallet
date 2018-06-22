package services

import (
	"entities"
	"sync"
	"utils"
	"dao"
	"rpcs"
)

type notifyService struct {
	BaseService
	sync.Once
	procsDeposits []entities.BaseDeposit
}

var _notifyService *notifyService

func GetNotifyService() *notifyService {
	if _notifyService == nil {
		_notifyService = new(notifyService)
		_notifyService.Once = sync.Once {}
		_notifyService.Once.Do(func() {
			_notifyService.create()
		})
	}
	return _notifyService
}

func (service *notifyService) create() error {
	service.status.RegAsObs(service)
	return service.BaseService.create()
}

func (service *notifyService) BeforeTurn(s *utils.Status, tgtStt int) {
	var err error
	switch tgtStt {
	case INIT:
		utils.LogMsgEx(utils.INFO, "initialization", nil)
		// 加载数据库中所有未稳定的充币交易
		if err = service.loadIncompleteDeposits(); err != nil {
			panic(utils.LogMsgEx(utils.ERROR, "加载未完成的充币失败：%v", err))
		}
	case START:
		utils.LogMsgEx(utils.INFO, "start", nil)
	}
}

func (service *notifyService) AfterTurn(s *utils.Status, srcStt int) {
	switch s.Current() {
	case INIT:
		utils.LogMsgEx(utils.INFO, "initialized", nil)
	case START:
		// 开启协程等待充币交易稳定
		go service.startWaitForStable()
		// 开启协程等待接收充币服务发来的交易
		go service.waitForUnstableDeposit()
		utils.LogMsgEx(utils.INFO, "started", nil)
	}
}

func (service *notifyService) loadIncompleteDeposits() error  {
	coinSetting := utils.GetConfig().GetCoinSettings()
	var err error
	service.procsDeposits, err = dao.GetDepositDAO().GetUnstableDeposit(coinSetting.Name)
	return err
}

func (service *notifyService) startWaitForStable() {
	var err error
	for err == nil && service.status.Current() == START {
		//如果达到协程上限，则等待
		coinSet := utils.GetConfig().GetCoinSettings()
		for _, deposit := range service.procsDeposits {
			// 获取当前块高
			var curHeight uint64
			if curHeight, err = rpcs.GetRPC(coinSet.Name).GetCurrentHeight(); err != nil {
				utils.LogMsgEx(utils.ERROR, "获取块高失败：%v", err)
				continue
			}

			stableHeight := uint64(coinSet.Stable)
			if deposit.Height + stableHeight >= curHeight {
				if err = TxIntoStable(&deposit, false); err != nil {
					continue
				}
			}
		}
	}
	service.status.TurnTo(DESTORY)
}

func (service *notifyService) waitForUnstableDeposit() {
	var err error
	for err == nil && service.status.Current() == START {
		var deposit entities.BaseDeposit
		var ok bool
		if deposit, ok = <- toNotifySig; !ok {
			break
		}
		utils.LogMsgEx(utils.INFO, "接收到一笔需等待的充币：%v", deposit)

		if _, err = dao.GetDepositDAO().AddScannedDeposit(&deposit); err != nil {
			utils.LogMsgEx(utils.ERROR, "添加未稳定提币记录失败：%v", err)
			continue
		}

		service.procsDeposits = append(service.procsDeposits, deposit)
		utils.LogMsgEx(utils.INFO, "充币交易（%s）已处于等待状态", deposit.TxHash)
	}
	service.status.TurnTo(DESTORY)
}