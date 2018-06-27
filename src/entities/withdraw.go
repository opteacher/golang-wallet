package entities

import (
	"time"
)

const (
	WITHDRAW_LOAD = 1
	WITHDRAW_SENT
	WITHDRAW_INCHAIN
	WITHDRAW_FINISHED
)

type BaseWithdraw struct {
	Address string
	Amount float64
	Asset string
}

type DatabaseWithdraw struct {
	BaseDeposit
	Id int
	TxHash string
	Status int
	Height uint64
	TxIndex int
	CreateTime time.Time
	UpdateTime time.Time
}

func TurnToBaseWithdraw(wd *DatabaseWithdraw) BaseWithdraw {
	var ret BaseWithdraw
	ret.Asset = wd.Asset
	ret.Address = wd.Address
	ret.Amount = wd.Amount
	return ret
}