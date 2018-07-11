package entities

import (
	"time"
)

const (
	DEPOSIT_FOUND = iota + 1
	DEPOSIT_FINISHED
)

type BaseDeposit struct {
	Transaction
	Address string	`json:"address"`
}

type DatabaseDeposit struct {
	BaseDeposit
	Id int					`json:"id"`
	Status int				`json:"status"`
	UpdateTime time.Time	`json:"update_time"`
}

func TurnTxToDeposit(tx *Transaction) BaseDeposit {
	var deposit BaseDeposit
	deposit.TxHash = tx.TxHash
	deposit.Address = tx.To
	deposit.Height = tx.Height
	deposit.CreateTime = tx.CreateTime
	deposit.Amount = tx.Amount
	deposit.Asset = tx.Asset
	deposit.TxIndex = tx.TxIndex
	return deposit
}