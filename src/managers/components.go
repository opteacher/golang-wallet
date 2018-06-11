package managers

import (
	"dao"
	"errors"
)

const (
	DEPOSIT_SERVICE = iota
	ADDRESS_DAO
)

type Component interface {
	Create() error
	IsCreate() bool
}

var components = map[int]Component {
	DEPOSIT_SERVICE:		NewDepositSvc(),
	ADDRESS_DAO:		dao.NewAddressDAO(),
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