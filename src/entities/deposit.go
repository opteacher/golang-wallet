package entities

import (
	"time"
)

const (
	DEPOSIT_FOUND = 1
	DEPOSIT_FINISHED
)

type BaseDeposit struct {
	TxHash string
	Address string
	Amount float64
	Asset string
	Height uint64
	TxIndex int
	CreateTime time.Time
}

type DatabaseDeposit struct {
	BaseDeposit
	Id int
	Status int
	UpdateTime time.Time
}