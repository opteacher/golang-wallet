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
	Address string
}

type DatabaseDeposit struct {
	BaseDeposit
	Id int
	Status int
	UpdateTime time.Time
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