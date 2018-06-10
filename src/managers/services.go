package managers

import (
	"utils"
	"log"
)

const (
	CREATE = iota
	INIT
	START
	PAUSE
	STOP
	DESTORY
)

type Service struct {
	status utils.Status
}

func (service *Service) BeforeTurn(s *utils.Status, tgtStt int) {
	switch tgtStt {
	case INIT:
		log.Println("initialization")
	case START:
		log.Println("start")
	}
}

func (service *Service) AfterTurn(s *utils.Status, srcStt int) {
	switch s.Current() {
	case INIT:
		log.Println("initialized")
	case START:
		log.Println("started")
	}
}

func (service *Service) Init() error {
	service.status.RegAsObs(service)
	service.status.Init([]int { CREATE, INIT, START })
	service.status.TurnTo(INIT)
	return nil
}

func (service *Service) Start() error {
	service.status.TurnTo(START)
	return nil
}