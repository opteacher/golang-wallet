package services

import (
	"utils"
	"log"
	"sync"
	"dao"
	"entities"
	"rpcs"
)

const (
	DESTORY = iota
	NONE
	CREATE
	INIT
	START
	PAUSE
	STOP
)

type DepositService struct {
	sync.Once
	status utils.Status
	addresses []string
	procsDeposits []entities.TotalDeposit
	height uint64
}

var _self *DepositService

func GetDepositService() *DepositService {
	if _self == nil {
		_self = new(DepositService)
		_self.Once = sync.Once {}
		_self.Once.Do(func() {
			_self.create()
		})
	}
	return _self
}

func (service *DepositService) create() error {
	service.status.RegAsObs(service)
	service.status.Init([]int { DESTORY, CREATE, INIT, START })
	return nil
}

func (service *DepositService) BeforeTurn(s *utils.Status, tgtStt int) {
	var err error
	switch tgtStt {
	case INIT:
		log.Println("initialization")
		// Load all address
		if err = service.loadAddresses(); err != nil {
			log.Fatal(err)
		}
		// Load all unstable deposits
		if err = service.loadIncompleteDeposits(); err != nil {
			log.Fatal(err)
		}
		// Get current height
		if err = service.getCurrentHeight(); err != nil {
			log.Fatal(err)
		}
	case START:
		log.Println("start")
	}
}

func (service *DepositService) AfterTurn(s *utils.Status, srcStt int) {
	var err error
	switch s.Current() {
	case INIT:
		log.Println("initialized")
	case START:
		log.Println("started")
		// Start goroutine to scan block chain
		if err = service.startScanChain(); err != nil {
			log.Fatal(err)
		}
	}
}

func (service *DepositService) Init() error {
	service.status.TurnTo(INIT)
	return nil
}

func (service *DepositService) Start() error {
	service.status.TurnTo(START)
	return nil
}

func (service *DepositService) loadAddresses() error {
	coinSetting := utils.GetConfig().GetCoinSettings()
	var err error
	service.addresses, err = dao.GetAddressDAO().FindInuseByAsset(coinSetting.Name)
	return err
}

func (service *DepositService) loadIncompleteDeposits() error  {
	coinSetting := utils.GetConfig().GetCoinSettings()
	var err error
	service.procsDeposits, err = dao.GetDepositDAO().GetUnstableDeposit(coinSetting.Name)
	return err
}

func (service *DepositService) getCurrentHeight() error {
	coinSetting := utils.GetConfig().GetCoinSettings()
	var err error
	service.height, err = dao.GetHeightDAO().GetHeight(coinSetting.Name)
	if service.height == 0 {
		_, err = dao.GetHeightDAO().ChkOrAddAsset(coinSetting.Name)
	}
	return err
}

func (service *DepositService) startScanChain() error {
	var err error
	coinName := utils.GetConfig().GetCoinSettings().Name
	rpc := rpcs.GetEth()
	depositDao := dao.GetDepositDAO()
	for ; err == nil && service.status.Current() == START; service.height++ {
		log.Printf("height: %d\n", service.height)

		// 获取当前块高
		var curHeight uint64
		if curHeight, err = rpc.GetCurrentHeight(); err != nil {
			log.Printf("Get current height failed: %s\n", err)
			continue
		}
		// 已经达到最高快高
		if service.height >= curHeight {
			continue
		}

		// 获取指定高度的交易
		var deposits []entities.BaseDeposit
		if deposits, err = rpc.GetTransactions(uint(service.height), service.addresses); err != nil {
			log.Printf("Get transaction failed: %s\n", err)
			continue
		}

		for _, deposit := range deposits {
			// 如果已经达到稳定块高，直接存入数据库
			// @tobo: 通知后台
			if deposit.Height + uint64(rpc.Stable) >= curHeight {
				if _, err = depositDao.AddStableDeposit(deposit); err != nil {
					log.Printf("Add deposit failed: %s\n", err)
					continue
				}
			} else {
				// 未进入稳定状态，抛给通知等待服务
			}
		}

		// 持久化高度到height表
		if service.height % 50 == 0 {
			if _, err = dao.GetHeightDAO().UpdateHeight(coinName, service.height); err != nil {
				log.Println("Update height failed: %s\n", err)
				continue
			}
		}
	}
	service.status.TurnTo(STOP)
	return err
}