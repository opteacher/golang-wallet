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
	Id int			`json:"id"`
	Address string	`json:"address"`
}

type DatabaseWithdraw struct {
	BaseWithdraw
	Status int				`json:"status"`
	UpdateTime time.Time	`json:"update_time"`
}

func TurnToBaseWithdraw(wd *DatabaseWithdraw) BaseWithdraw {
	var ret BaseWithdraw
	ret.Asset = wd.Asset
	ret.Address = wd.Address
	ret.Amount = wd.Amount
	return ret
}