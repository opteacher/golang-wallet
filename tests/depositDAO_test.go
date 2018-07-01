package tests

import (
	"dao"
	"entities"
	"testing"
	"fmt"
)

func TestAddScannedDepositDAO(t *testing.T) {
	var deposit entities.BaseDeposit
	deposit.TxHash	= "0x12345"
	deposit.Address	= "0xabcd"
	deposit.Amount	= 1000
	deposit.TxIndex	= 0
	deposit.Height	= 200000
	deposit.Asset	= "ETH"
	fmt.Println(dao.GetDepositDAO().AddScannedDeposit(&deposit))
}

func TestAddStableDepositDAO(t *testing.T) {
	var deposit entities.BaseDeposit
	deposit.TxHash	= "0x12345"
	deposit.Address	= "0xabcd"
	deposit.Amount	= 1000
	deposit.TxIndex	= 0
	deposit.Height	= 200000
	deposit.Asset	= "ETH"
	fmt.Println(dao.GetDepositDAO().AddStableDeposit(&deposit))
}

func TestGetUnstableDepositDAO(t *testing.T) {
	fmt.Println(dao.GetDepositDAO().GetUnstableDeposit("ETH"))
}

func TestDepositIntoStableDAO(t *testing.T) {
	fmt.Println(dao.GetDepositDAO().DepositIntoStable("0x12345"))
}