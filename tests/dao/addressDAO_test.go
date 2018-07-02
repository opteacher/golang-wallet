package dao

import (
	"dao"
	"testing"
	"fmt"
)

func TestNewAddressDAO(t *testing.T) {
	fmt.Println(dao.GetAddressDAO().NewAddress("BTC", "0xabcd"))
}

func TestFindAddressByAsset(t *testing.T) {
	fmt.Println(dao.GetAddressDAO().FindInuseByAsset("ETH"))
}
