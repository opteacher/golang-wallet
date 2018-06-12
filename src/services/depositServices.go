package services

import (
	"utils"
	"log"
	"sync"
	"dao"
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
}

func (service *DepositService) create() error {
	service.status.RegAsObs(service)
	service.status.Init([]int { DESTORY, CREATE, INIT, START })
	return nil
}

func (service *DepositService) IsCreate() bool {
	return service.status.Current() > CREATE
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
	case START:
		log.Println("start")
	}
}

func (service *DepositService) AfterTurn(s *utils.Status, srcStt int) {
	switch s.Current() {
	case INIT:
		log.Println("initialized")
	case START:
		log.Println("started")
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