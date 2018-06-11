package dao

import (
	"log"
	"errors"
	"databases"
	"database/sql"
)

const TABLE_NAME = "address"
const CREATE_SQL = `CREATE TABLE IF NOT EXISTS address (
	id INTEGER NOT NULL AUTO_INCREMENT,
	asset VARCHAR(255) NOT NULL,
	address VARCHAR(255) NOT NULL UNIQUE,
	inuse TINYINT(1) NOT NULL DEFAULT 0
) ENGINE=InnoDB DEFAULT CHARSET=utf8`

type AddressDao struct {
	created bool
}

func (dao *AddressDao) Create() error {
	var db *sql.DB
	var err error
	if db, err = databases.Connect(); err != nil {
		log.Println(err)
		return err
	}
	if _, err = db.Exec(CREATE_SQL); err != nil {
		log.Println(err)
		return err
	}
	dao.created = true
	return nil
}

func (dao *AddressDao) IsCreate() bool {
	return dao.created
}

func (dao *AddressDao) FindByAsset(asset string) ([]string, error) {
	if !dao.created {
		return nil, errors.New("Hasnt created")
	}
	var db *sql.DB
	var err error
	if db, err = databases.Connect(); err != nil {
		log.Println(err)
		return nil, err
	}
	SQL := "SELECT address FROM address WHERE inuse=1 AND asset=?"
	var rows *sql.Rows
	if rows, err = db.Query(SQL, asset); err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()
	var addresses []string
	for rows.Next() {
		var address string
		if err = rows.Scan(&address); err != nil {
			log.Println(err)
			continue
		}
		addresses = append(addresses, address)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
	return addresses, nil
}

func NewAddressDAO() *AddressDao {
	ret := new(AddressDao)
	ret.created = false
	return ret
}