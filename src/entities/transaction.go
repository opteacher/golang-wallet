package entities

import "time"

type Transaction struct {
	TxHash string			`json:"tx_hash"`
	BlockHash string		`json:"block_hash"`
	From string				`json:"from"`
	To string				`json:"to"`
	Amount float64			`json:"amount"`
	Asset string			`json:"asset"`
	Height uint64			`json:"height"`
	TxIndex int				`json:"tx_index"`
	CreateTime time.Time	`json:"create_time"`
}