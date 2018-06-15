package services

import (
	"utils"
	"sync"
	"dao"
	"entities"
	"rpcs"
)

type depositService struct {
	sync.Once
	status utils.Status
	addresses []string
	height uint64
}

var _depositService *depositService

func GetDepositService() *depositService {
	if _depositService == nil {
		_depositService = new(depositService)
		_depositService.Once = sync.Once {}
		_depositService.Once.Do(func() {
			_depositService.create()
		})
	}
	return _depositService
}

func (service *depositService) create() error {
	service.status.RegAsObs(service)
	service.status.Init([]int { DESTORY, CREATE, INIT, START, STOP })
	return nil
}

func (service *depositService) BeforeTurn(s *utils.Status, tgtStt int) {
	var err error
	switch tgtStt {
	case INIT:
		utils.LogMsgEx(utils.INFO, "initialization", nil)
		// 加载所有内部地址
		if err = service.loadAddresses(); err != nil {
			panic(utils.LogMsgEx(utils.ERROR, "加载地址失败：%v", err))
		}
		// 获取上一次扫描的块高
		if err = service.getCurrentHeight(); err != nil {
			panic(utils.LogMsgEx(utils.ERROR, "获取当前块高失败：%v", err))
		}
	case START:
		utils.LogMsgEx(utils.INFO, "start", nil)
	}
}

func (service *depositService) AfterTurn(s *utils.Status, srcStt int) {
	switch s.Current() {
	case INIT:
		utils.LogMsgEx(utils.INFO, "initialized", nil)
	case START:
		// 开启协程扫描区块链上的交易记录
		go service.startScanChain()
		utils.LogMsgEx(utils.INFO, "started", nil)
	}
}

func (service *depositService) Init() {
	service.status.TurnTo(INIT)
}

func (service *depositService) Start() {
	service.status.TurnTo(START)
}

func (service *depositService) Stop()  {
	service.status.TurnTo(STOP)
}

func (service *depositService) IsDestroy() bool {
	return service.status.Current() == DESTORY
}

func (service *depositService) loadAddresses() error {
	coinSetting := utils.GetConfig().GetCoinSettings()
	var err error
	service.addresses, err = dao.GetAddressDAO().FindInuseByAsset(coinSetting.Name)
	return err
}

func (service *depositService) getCurrentHeight() error {
	coinSetting := utils.GetConfig().GetCoinSettings()
	var err error
	service.height, err = dao.GetHeightDAO().GetHeight(coinSetting.Name)
	if service.height == 0 {
		_, err = dao.GetHeightDAO().ChkOrAddAsset(coinSetting.Name)
	}
	return err
}

func (service *depositService) startScanChain() {
	var err error
	coinCfg := utils.GetConfig().GetCoinSettings()
	rpc := rpcs.GetEth()
	for err == nil && service.status.Current() == START {
		utils.LogMsgEx(utils.INFO, "块高: %d", service.height)

		// 获取指定高度的交易
		var deposits []entities.BaseDeposit
		if deposits, err = rpc.GetTransactions(uint(service.height), service.addresses); err != nil {
			utils.LogMsgEx(utils.ERROR, "获取交易失败：%v", err)
			continue
		}

		for _, deposit := range deposits {
			utils.LogMsgEx(utils.INFO, "发现交易：%v", deposit)
			dao.GetProcessDAO().SaveProcess(&entities.DatabaseProcess {
				entities.BaseProcess {
					deposit.TxHash,
					entities.DEPOSIT,
					entities.NOTIFY,
					true,
				},
				deposit.Height,
				deposit.Height + uint64(coinCfg.Stable),
			})

			// 获取当前块高
			var curHeight uint64
			if curHeight, err = rpc.GetCurrentHeight(); err != nil {
				utils.LogMsgEx(utils.ERROR, "获取块高失败：%v", err)
				continue
			}

			// 如果已经达到稳定块高，直接存入数据库
			// @tobo: 通知后台
			if deposit.Height + uint64(rpc.Stable) >= curHeight {
				utils.LogMsgEx(utils.INFO, "交易（%s）进入稳定状态", deposit.TxHash)

				if err = TxIntoStable(&deposit, true); err != nil {
					utils.LogMsgEx(utils.ERROR, "插入稳定交易记录失败：%v", err)
					continue
				}
			} else {
				// 未进入稳定状态，抛给通知等待服务
				toNotifySig <- deposit
				utils.LogMsgEx(utils.INFO, "交易（%s）进入等待列队", deposit.TxHash)
			}
		}

		// 持久化高度到height表
		if service.height % 20 == 0 {
			if _, err = dao.GetHeightDAO().UpdateHeight(coinCfg.Name, service.height); err != nil {
				utils.LogMsgEx(utils.ERROR, "更新块高失败：%v", err)
				continue
			}
			utils.LogMsgEx(utils.INFO, "块高更新到：%d", service.height)
		}

		service.height++
	}
	close(toNotifySig)
	service.status.TurnTo(DESTORY)
}