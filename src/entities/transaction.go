package entities

import "time"

type Transaction struct {
	TxHash string
	From string
	To string
	Amount float64
	Asset string
	Height uint64
	TxIndex int
	CreateTime time.Time
}