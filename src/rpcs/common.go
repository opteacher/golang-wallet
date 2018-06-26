package rpcs

import (
	"entities"
	"reflect"
)

type Rpc interface {
	GetTransactions(height uint, addresses []string) ([]entities.BaseDeposit, error)
	GetCurrentHeight() (uint64, error)
	GetDepositAmount() (map[string]float64, error)
	GetBalance(address string) (float64, error)
	SendFrom(from string, to string, amount float64) (string, error)
	GetNewAddress() (string, error)
}

type rpc struct {

}

var __rpc = new(rpc)

func GetRPC(name string) Rpc {
	return reflect.ValueOf(__rpc).MethodByName(name).Call(nil)[0].Interface().(Rpc)
}

type RequestBody struct {
	Method string			`json:method`
	Params []interface{}	`json:params`
	Id string				`json:id`
}
