package services

import (
	"entities"
	"sync"
	"utils"
	"dao"
	"rpcs"
)

/***
	通知服务：接收来自充值和提币的交易，等待其稳定后做后续操作
	子协程（startWaitForStable）：等待交易进入稳定状态
	子协程（waitForUnstableTransaction）：等待来自充值和提币的交易
 */
type notifyService struct {
	BaseService
	sync.Once
	waitForStableTxs []entities.Transaction
	waitTxsCounter map[string]uint
	waitForStableTxsLock *sync.Mutex
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
	service.name = "notifyService"
	service.status.RegAsObs(service)
	service.waitTxsCounter = make(map[string]uint)
	service.waitForStableTxsLock = new(sync.Mutex)
	return service.BaseService.create()
}

func (service *notifyService) BeforeTurn(s *utils.Status, tgtStt int) {
	var err error
	switch tgtStt {
	case INIT:
		utils.LogMsgEx(utils.INFO, "initialization", nil)
		// 加载数据库中所有未稳定的充币交易
		if err = service.loadIncompleteTransactions(); err != nil {
			panic(utils.LogMsgEx(utils.ERROR, "加载未完成的交易失败：%v", err))
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
		go service.waitForUnstableTransaction()
		utils.LogMsgEx(utils.INFO, "started", nil)
	}
}

func (service *notifyService) loadIncompleteTransactions() error  {
	coinSetting := utils.GetConfig().GetCoinSettings()
	var err error
	var deposits []entities.BaseDeposit
	if deposits, err = dao.GetDepositDAO().GetUnstableDeposit(coinSetting.Name); err != nil {
		return err
	}
	var withdraws []entities.DatabaseWithdraw
	if withdraws, err = dao.GetWithdrawDAO().GetUnstableWithdraw(coinSetting.Name); err != nil {
		return err
	}
	for _, deposit := range deposits {
		utils.LogMsgEx(utils.INFO, "发现一笔待稳定的充值记录：%s", deposit.TxHash)
		service.waitForStableTxs = append(service.waitForStableTxs, deposit.Transaction)
		service.waitTxsCounter[deposit.TxHash] = 0
	}
	for _, withdraw := range withdraws {
		utils.LogMsgEx(utils.INFO, "发现一笔待稳定的提币记录：%s", withdraw.TxHash)
		service.waitForStableTxs = append(service.waitForStableTxs, withdraw.Transaction)
		service.waitTxsCounter[withdraw.TxHash] = 0
	}
	return err
}

func (service *notifyService) startWaitForStable() {
	var err error
	coinSet := utils.GetConfig().GetCoinSettings()
	for err == nil && service.status.Current() == START {
		//如果达到协程上限，则等待
		service.waitForStableTxsLock.Lock()
		for i, tx := range service.waitForStableTxs {
			// 获取当前块高
			var curHeight uint64
			if curHeight, err = rpcs.GetRPC(coinSet.Name).GetCurrentHeight(); err != nil {
				utils.LogMsgEx(utils.ERROR, "获取块高失败：%v", err)
				continue
			}

			stableHeight := uint64(coinSet.Stable)
			if curHeight >= tx.Height + stableHeight {
				utils.LogMsgEx(utils.INFO, "交易：%s已进入稳定状态", tx.TxHash)

				if err = TxIntoStable(tx.TxHash, tx.Asset); err != nil {
					continue
				}

				service.waitForStableTxs = append(service.waitForStableTxs[:i], service.waitForStableTxs[i + 1:]...)
				delete(service.waitTxsCounter, tx.TxHash)
				break
			} else {
				if service.waitTxsCounter[tx.TxHash] % 20 == 0 {
					utils.LogMsgEx(utils.INFO, "交易：%s等待稳定，%d/%d",
						tx.TxHash, curHeight, tx.Height + stableHeight)
				}
				service.waitTxsCounter[tx.TxHash]++
			}
		}
		service.waitForStableTxsLock.Unlock()
	}
	service.status.TurnTo(DESTORY)
}

func (service *notifyService) waitForUnstableTransaction() {
	var err error
	for err == nil && service.status.Current() == START {
		var tx entities.Transaction
		var ok bool
		if tx, ok = <- toNotifySig; !ok {
			break
		}
		utils.LogMsgEx(utils.INFO, "接收到一笔需等待的交易：%v", tx)

		service.waitForStableTxsLock.Lock()
		service.waitForStableTxs = append(service.waitForStableTxs, tx)
		service.waitTxsCounter[tx.TxHash] = 0
		service.waitForStableTxsLock.Unlock()
		utils.LogMsgEx(utils.INFO, "交易（%s）已处于等待状态", tx.TxHash)
	}
	service.status.TurnTo(DESTORY)
}