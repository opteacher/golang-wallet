package entities

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
	CONFIRM	= "CONFIRM"
	NOTIFY	= "NOTIFY"
	FINISH	= "FINISH"
)

var Processes = []string {
	AUDIT, LOAD, SENT, SENDING, CONFIRM, NOTIFY, FINISH,
}

type BaseProcess struct {
	TxHash string		`field:tx_hash`
	Type string			`field:type`
	Process string		`field:process`
	Cancelable bool		`field:cancelable`
}

type DatabaseProcess struct {
	BaseProcess
	Height uint64		`field:height`
	CompleteHeight uint64	`field:complete_height`
}