package tests

import (
	"log"
	"rpcs"
	"testing"
	"entities"
	"dao"
)

func TestGetTransactions(t *testing.T) {
	log.SetFlags(log.Lshortfile)

	//Test RPC
	var err error
	var totalAffectRows int64
	var tempAffectRows int64
	var txs []entities.BaseDeposit
	txs, err = rpcs.GetRPC("ETH").GetTransactions(120, []string {
		"0x43faead79328ca23fbb179af73ab8c0153ed990c",
	})
	totalAffectRows = 0
	depositDAO := dao.GetDepositDAO()
	for _, tx := range txs {
		if tempAffectRows, err = depositDAO.AddScannedDeposit(&tx); err != nil {
			log.Fatal(err)
		}
		totalAffectRows += tempAffectRows
	}
	log.Printf("Add deposits succeed: %d\n", totalAffectRows)
}

func TestSendTransaction(t *testing.T) {
	var err error
	var txHash string
	if txHash, err = rpcs.GetRPC("ETH").SendFrom(
		"0x47e8e8e49a8c1c308e84439f2c55ef18710f5ed6",
		"0x65711d2e616b437e65273d30d4385fd0028a461b",
		10.23456); err != nil {
		log.Fatal(err)
	}
	log.Println(txHash)
}

func TestNewAddress(t *testing.T) {
	var addr string
	var err error
	if addr, err = rpcs.GetRPC("ETH").GetNewAddress(); err != nil {
		log.Fatal(err)
	} else {
		log.Println(addr)
	}
}