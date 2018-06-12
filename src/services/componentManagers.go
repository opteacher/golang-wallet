package services

import (
	"dao"
	"errors"
	"utils"
)

const (
	DEPOSIT_SERVICE = iota
	ADDRESS_DAO
	CONFIG
)

var components = map[int]Component {
	DEPOSIT_SERVICE:	NewDepositSvc(),
	ADDRESS_DAO:		dao.NewAddressDAO(),
	CONFIG:				utils.NewConfig(),
}

func GetComponent(name int) (Component, error) {
	component, ok := components[name]
	if !ok {
		return nil, errors.New("Could not find identitified component")
	}
	if !component.IsCreate() {
		if err := component.Create(); err != nil {
			return component, err
		}
	}
	return component, nil
}