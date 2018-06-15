package main

import (
	"dao"
	"log"
)

func main() {
	log.SetFlags(log.Lshortfile)

	// Test DAO
	addressDAO := dao.GetAddressDAO()
	addressDAO.NewAddress("ETH", "0xabcd")
	addressDAO.NewAddressInuse("BTC", "0x1234")
	log.Println(addressDAO.FindInuseByAsset("BTC"))
}
