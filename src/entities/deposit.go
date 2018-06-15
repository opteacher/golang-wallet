package entities

import (
	"time"
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