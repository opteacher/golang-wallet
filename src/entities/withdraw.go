package entities

import (
	"time"
)

const (
	WITHDRAW_LOAD = iota + 1
	WITHDRAW_SENT
	WITHDRAW_INCHAIN
	WITHDRAW_FINISHED
)

type BaseWithdraw struct {
	Transaction
	Id int
	Address string
}

type DatabaseWithdraw struct {
	BaseDeposit
	Status int
	UpdateTime time.Time
}

func TurnToBaseWithdraw(wd *DatabaseWithdraw) BaseWithdraw {
	var ret BaseWithdraw
	ret.Asset = wd.Asset
	ret.Address = wd.Address
	ret.Amount = wd.Amount
	return ret
}