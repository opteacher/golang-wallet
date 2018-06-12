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
	inuse TINYINT(1) NOT NULL DEFAULT 0,
	PRIMARY KEY(id)
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

func (dao *AddressDao) NewAddress(asset string, address string) (int64, error) {
	var db *sql.DB
	var err error
	if db, err = dao.connect(); err != nil {
		log.Println(err)
		return 0, err
	}
	SQL := "INSERT INTO address (asset, address) VALUES (?, ?)"
	var result sql.Result
	if result, err = db.Exec(SQL, asset, address); err != nil {
		log.Println(err)
		return 0, err
	}
	return result.RowsAffected()
}

func (dao *AddressDao) FindByAsset(asset string) ([]string, error) {
	var db *sql.DB
	var err error
	if db, err = dao.connect(); err != nil {
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

func (dao *AddressDao) connect() (*sql.DB, error) {
	if !dao.created {
		return nil, errors.New("Hasnt created")
	}
	return databases.Connect()
}

func NewAddressDAO() *AddressDao {
	ret := new(AddressDao)
	ret.created = false
	return ret
}