package main

import (
	"entities"
	"rpcs"
	"log"
	"dao"
)

func main() {
	log.SetFlags(log.Lshortfile)

	//Test RPC
	var err error
	var totalAffectRows int64
	var tempAffectRows int64
	var txs []entities.BaseDeposit
	txs, err = rpcs.GetEth().GetTransactions(120, []string {
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