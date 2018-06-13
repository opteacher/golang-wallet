package entities

import (
	"time"
)

type Deposit struct {
	TxHash string
	Address string
	Amount float64
	Asset string
	Height uint64
	TxIndex int
	CreateTime time.Time
}