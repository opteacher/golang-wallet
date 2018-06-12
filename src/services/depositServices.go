package services

import (
	"utils"
	"log"
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
	status utils.Status
	addresses []string
}

func (service *DepositService) Create() error {
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
		var componentTmp Component
		if componentTmp, err = GetComponent(ADDRESS_DAO); err != nil {
			log.Fatal(err)
		}
		addressDAO := componentTmp.(*dao.AddressDao)
		if service.addresses, err = addressDAO.FindByAsset("ETH"); err != nil {
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
	return nil
}

func NewDepositSvc() *DepositService {
	return new(DepositService)
}