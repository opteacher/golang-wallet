package dao

import (
	"testing"
	"entities"
	"fmt"
	"dao"
)

func TestNewWithdraw(t *testing.T) {
	var wd entities.BaseWithdraw
	wd.Id = 12345
	wd.TxHash = "0x12345"
	wd.Asset = "ETH"
	wd.Address = "0xabcd"
	wd.Amount = 123.5
	fmt.Println(dao.GetWithdrawDAO().NewWithdraw(wd))
}

func TestWithdrawIntoChain(t *testing.T) {
	fmt.Println(dao.GetWithdrawDAO().WithdrawIntoChain("0x12345", 12345, 0))
}

func TestWithdrawIntoStable(t *testing.T) {
	fmt.Println(dao.GetWithdrawDAO().WithdrawIntoStable("0x12345"))
}

func TestGetAllUnstable(t *testing.T) {
	fmt.Println(dao.GetWithdrawDAO().GetAllUnstable("ETH"))
}

func TestGetAvailableId(t *testing.T) {
	fmt.Println(dao.GetWithdrawDAO().GetAvailableId("ETH"))
}