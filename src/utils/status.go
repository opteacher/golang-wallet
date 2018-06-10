package utils

type Observer interface {
	BeforeTurn(status *Status)
	AfterTurn(status *Status)
}

type Status struct {
	AllStatus	[]int
	StatusVal	int
	Observers	[]Observer
}

func (stt *Status) Init(stts []int) {
	if stts != nil {
		stt.AllStatus = stts
	}
}

func (stt *Status) TurnTo(status int) int {
	orgStt := stt.StatusVal
	if !Contains(stt.AllStatus, status) {
		return orgStt
	}
	for _, obs := range stt.Observers {
		obs.BeforeTurn(stt)
	}
	stt.StatusVal = status
	for _, obs := range stt.Observers {
		obs.AfterTurn(stt)
	}
	return orgStt
}

func (stt *Status) Current() int {
	return stt.StatusVal
}