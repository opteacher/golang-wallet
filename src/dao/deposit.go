package dao

import (
	"sync"
	"database/sql"
	"entities"
	"unsafe"
	"time"
	"strconv"
	"utils"
	"github.com/go-sql-driver/mysql"
)

type depositDao struct {
	baseDao
	sync.Once
}

var _depositDao *depositDao

func GetDepositDAO() *depositDao {
	if _depositDao == nil {
		_depositDao = new(depositDao)
		_depositDao.Once = sync.Once {}
		_depositDao.Once.Do(func() {
			_depositDao.create("deposit")
		})
	}
	return _depositDao
}

func (d *depositDao) AddScannedDeposit(deposit *entities.BaseDeposit) (int64, error) {
	var params = []interface {} {
		deposit.TxHash,
		deposit.Address,
		deposit.Amount,
		deposit.Asset,
		deposit.Height,
		deposit.TxIndex,
	}
	var useSQL = "AddScannedDeposit"
	if deposit.CreateTime.Year() > 1000 {
		useSQL = "AddDepositWithTime"
		params = append(params, deposit.CreateTime)
	}
	return insertTemplate((*baseDao)(unsafe.Pointer(d)), useSQL, params)
}

func (d *depositDao) AddStableDeposit(deposit *entities.BaseDeposit) (int64, error) {
	if deposit.CreateTime.Year() < 1000 {
		deposit.CreateTime = time.Now()
	}
	return insertTemplate((*baseDao)(unsafe.Pointer(d)), "AddStableDeposit", []interface {} {
		deposit.TxHash,
		deposit.Address,
		deposit.Amount,
		deposit.Asset,
		deposit.Height,
		deposit.TxIndex,
		deposit.CreateTime,
	})
}

func (d *depositDao) GetUnstableDeposit(asset string) ([]entities.BaseDeposit, error) {
	bd := (*baseDao)(unsafe.Pointer(d))
	conds := []interface {} { asset }
	var result []map[string]interface {}
	var err error
	if result, err = selectTemplate(bd, "GetUnstableDeposit", conds); err != nil {
		return nil, err
	}

	var ret []entities.BaseDeposit
	for _, entity := range result {
		var deposit entities.BaseDeposit
		deposit.To = string(*entity["address"].(*sql.RawBytes))
		deposit.Address = deposit.To
		deposit.Amount, err = strconv.ParseFloat(string(*entity["amount"].(*sql.RawBytes)), 64)
		if err != nil {
			panic(utils.LogMsgEx(utils.ERROR, "解析交易金额失败：%v", err))
		}
		deposit.Asset = asset
		deposit.TxHash = string(*entity["tx_hash"].(*sql.RawBytes))
		deposit.Height = uint64(*entity["height"].(*int32))
		deposit.TxIndex = int(entity["tx_index"].(*sql.NullInt64).Int64)
		ret = append(ret, deposit)
	}
	return ret, nil
}

func (d *depositDao) DepositIntoStable(txHash string) (int64, error) {
	return updatePartsTemplate((*baseDao)(unsafe.Pointer(d)), "DepositIntoStable",
		[]interface {} { txHash }, nil)
}

func (d *depositDao) GetDepositId(txHash string) (int, error) {
	var result []map[string]interface {}
	var err error
	bd := (*baseDao)(unsafe.Pointer(d))
	conds := []interface {} { txHash }
	if result, err = selectTemplate(bd, "GetDepositId", conds); err != nil {
		return -1, err
	}

	if len(result) != 1 {
		return -1, utils.LogMsgEx(utils.ERROR, "找不到交易：%s的充币ID", txHash)
	}
	depositId := (int) (*(result[0]["id"].(*int32)))
	return depositId, nil
}

func (d *depositDao)GetDeposits(conds map[string]interface {}) ([]entities.DatabaseDeposit, error) {
	bd := (*baseDao)(unsafe.Pointer(d))
	var result []map[string]interface {}
	var err error
	if result, err = selectPartsTemplate(bd, "GetDeposits", conds); err != nil {
		return nil, err
	}

	var ret []entities.DatabaseDeposit
	for _, entity := range result {
		var deposit entities.DatabaseDeposit
		deposit.Id = int(*entity["id"].(*int32))
		deposit.TxHash = string(*entity["tx_hash"].(*sql.RawBytes))
		deposit.To = string(*entity["address"].(*sql.RawBytes))
		deposit.Address = deposit.To
		deposit.Amount, err = strconv.ParseFloat(string(*entity["amount"].(*sql.RawBytes)), 64)
		if err != nil {
			panic(utils.LogMsgEx(utils.ERROR, "解析交易金额失败：%v", err))
		}
		deposit.Asset = string(*entity["asset"].(*sql.RawBytes))
		deposit.Height = uint64(*entity["height"].(*int32))
		deposit.TxIndex = int(entity["tx_index"].(*sql.NullInt64).Int64)
		deposit.CreateTime = entity["create_time"].(*mysql.NullTime).Time
		deposit.UpdateTime = entity["update_time"].(*mysql.NullTime).Time
		ret = append(ret, deposit)
	}
	return ret, nil
}