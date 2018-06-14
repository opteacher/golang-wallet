package services

import (
	"log"
	"dao"
	"entities"
)

func TxIntoStable(deposit *entities.BaseDeposit) error {
	var err error
	if _, err = dao.GetDepositDAO().AddStableDeposit(deposit); err != nil {
		log.Printf("Add deposit failed: %s\n", err)
		return err
	}
	return nil
}