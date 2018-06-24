package tests

import (
	"dao"
	"entities"
	"log"
	"testing"
)

func TestDepositDAO(t *testing.T) {
	log.SetFlags(log.Lshortfile)

	var err error
	var totalAffectRows int64

	depositDAO := dao.GetDepositDAO()
	var deposit entities.BaseDeposit
	deposit.TxHash	= "0x12345"
	deposit.Address	= "0xabcd"
	deposit.Amount	= 1000
	deposit.TxIndex	= 0
	deposit.Height	= 200000
	deposit.Asset	= "ETH"
	if totalAffectRows, err = depositDAO.AddScannedDeposit(&deposit); err != nil {
		log.Fatal(err)
	}
	log.Printf("Add deposit succeed: %d\n", totalAffectRows)
}