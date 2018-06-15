package main

import (
	"entities"
	"dao"
	"log"
)

func main() {
	log.SetFlags(log.Lshortfile)

	var err error
	var totalAffectRows int64

	service := dao.GetProcessDAO()

	var processForAdd entities.DatabaseProcess
	processForAdd.TxHash = "0xabcd"
	processForAdd.Cancelable = true
	processForAdd.Process = entities.AUDIT
	processForAdd.Type = entities.WITHDRAW
	if totalAffectRows, err = service.SaveProcess(&processForAdd); err != nil {
		log.Fatal(err)
	}
	log.Println(totalAffectRows)

	var processForUpd entities.DatabaseProcess
	processForUpd.TxHash = "0xabcd"
	processForUpd.Height = 2000
	processForUpd.Cancelable = true
	if totalAffectRows, err = service.SaveProcess(&processForUpd); err != nil {
		log.Fatal(err)
	}
	log.Println(totalAffectRows)
}
