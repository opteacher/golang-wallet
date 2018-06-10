package utils

import "errors"

type Observer interface {
	BeforeTurn(status *Status, tgtStt int)
	AfterTurn(status *Status, srcStt int)
}

type Status struct {
	AllStatus	[]int
	statusVal	int
	observers	[]Observer
}

func (stt *Status) Init(stts []int) {
	if stts != nil {
		stt.AllStatus = stts
	}
}

func (stt *Status) TurnTo(status int) (int, error) {
	orgStt := stt.statusVal
	if !Contains(stt.AllStatus, status) {
		return orgStt, errors.New("Could not find identified status")
	}
	for _, obs := range stt.observers {
		obs.BeforeTurn(stt, status)
	}
	stt.statusVal = status
	for _, obs := range stt.observers {
		obs.AfterTurn(stt, orgStt)
	}
	return orgStt, nil
}

func (stt *Status) Current() int {
	return stt.statusVal
}

func (stt *Status) RegAsObs(obs Observer) error {
	stt.observers = append(stt.observers, obs)
	return nil
}