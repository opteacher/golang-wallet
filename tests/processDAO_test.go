package tests

import (
	"entities"
	"dao"
	"testing"
	"fmt"
)

func TestInsertProcsDAO(t *testing.T) {
	var processForAdd entities.DatabaseProcess
	processForAdd.TxHash = "0xabcd"
	processForAdd.Cancelable = true
	processForAdd.Process = entities.AUDIT
	processForAdd.Type = entities.WITHDRAW
	processForAdd.Asset = "ETH"
	fmt.Println(dao.GetProcessDAO().SaveProcess(&processForAdd))
}

func TestUpdateProcsDAO(t *testing.T) {
	var processForUpd entities.DatabaseProcess
	processForUpd.TxHash = "0xabcd"
	processForUpd.Height = 2000
	processForUpd.Cancelable = true
	fmt.Println(dao.GetProcessDAO().SaveProcess(&processForUpd))
}

func TestQueryProcsDAO(t *testing.T) {
	fmt.Println(dao.GetProcessDAO().QueryProcess("ETH", "0xabcd"))
}