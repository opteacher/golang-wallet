package entities

import "time"

const (
	DEPOSIT		= "DEPOSIT"
	COLLECT		= "COLLECT"
	WITHDRAW	= "WITHDRAW"
)

var Types = []string {
	DEPOSIT, COLLECT, WITHDRAW,
}

const (
	AUDIT	= "AUDIT"
	LOAD	= "LOAD"
	SENT	= "SENT"
	SENDING	= "SENDING"
	INCHAIN	= "INCHAIN"
	FINISH	= "FINISH"
)

var Processes = []string {
	AUDIT, LOAD, SENT, SENDING, INCHAIN, FINISH,
}

type BaseProcess struct {
	Id int				`json:"id"`
	TxHash string		`json:"tx_hash"`
	Asset string		`json:"asset"`
	Type string			`json:"type"`
	Process string		`json:"process"`
	Cancelable bool		`json:"cancelable"`
}

type DatabaseProcess struct {
	BaseProcess
	Height uint64			`json:"height"`
	CurrentHeight uint64	`json:"current_height"`
	CompleteHeight uint64	`json:"complete_height"`
	LastUpdateTime time.Time`json:"last_update_time"`
}