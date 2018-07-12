package dao

import (
	"sync"
	"entities"
	"unsafe"
	"utils"
	"errors"
	"database/sql"
	"strconv"
	"github.com/go-sql-driver/mysql"
)

type withdrawDao struct {
	baseDao
	sync.Once
}

var _withdrawDao *withdrawDao

func GetWithdrawDAO() *withdrawDao {
	if _withdrawDao == nil {
		_withdrawDao = new(withdrawDao)
		_withdrawDao.Once = sync.Once {}
		_withdrawDao.Once.Do(func() {
			_withdrawDao.create("withdraw")
		})
	}
	return _withdrawDao
}

func (d *withdrawDao) getWithdraws(asset string, sqlName string) ([]entities.DatabaseWithdraw, error) {
	var result []map[string]interface {}
	var err error
	conds := []interface {} { asset }
	if result, err = selectTemplate((*baseDao)(unsafe.Pointer(d)), sqlName, conds); err != nil {
		return nil, utils.LogIdxEx(utils.ERROR, 39, err)
	}

	var ret []entities.DatabaseWithdraw
	for _, entity := range result {
		var bwd entities.DatabaseWithdraw
		bwd.Id = int(*entity["id"].(*int32))
		bwd.TxHash = string(*entity["tx_hash"].(*sql.RawBytes))
		bwd.Address = string(*entity["address"].(*sql.RawBytes))
		bwd.Amount, err = strconv.ParseFloat(string(*entity["amount"].(*sql.RawBytes)), 64)
		if err != nil {
			panic(utils.LogMsgEx(utils.ERROR, "解析交易金额失败：%v", err))
		}
		bwd.Asset = string(*entity["asset"].(*sql.RawBytes))
		bwd.Height = uint64(entity["height"].(*sql.NullInt64).Int64)
		bwd.TxIndex = int(entity["tx_index"].(*sql.NullInt64).Int64)
		bwd.Status = int(entity["status"].(*sql.NullInt64).Int64)
		bwd.CreateTime = entity["create_time"].(*mysql.NullTime).Time
		bwd.UpdateTime = entity["update_time"].(*mysql.NullTime).Time
		ret = append(ret, bwd)
	}
	return ret, nil
}

func (d *withdrawDao) GetUnfinishWithdraw(asset string) ([]entities.DatabaseWithdraw, error) {
	return d.getWithdraws(asset, "GetUnfinishWithdraw")
}

func (d *withdrawDao) GetUnstableWithdraw(asset string) ([]entities.DatabaseWithdraw, error) {
	return d.getWithdraws(asset, "GetUnstableWithdraw")
}

func (d *withdrawDao) GetAvailableId(asset string) (int, error) {
	var result []map[string]interface {}
	var err error
	conds := []interface {} { asset }
	if result, err = selectTemplate((*baseDao)(unsafe.Pointer(d)), "GetAvailableId", conds); err != nil {
		return 0, utils.LogIdxEx(utils.ERROR, 38, err)
	}
	if len(result) != 1 {
		return 0, utils.LogIdxEx(utils.ERROR, 38, errors.New("返回的id数量不等于1"))
	}

	newId := result[0]
	var ret *sql.NullInt64
	var ok bool
	if ret, ok = newId["new_id"].(*sql.NullInt64); !ok {
		return 0, utils.LogIdxEx(utils.ERROR, 38, errors.New("返回值不包含new_id"))
	}
	return int(ret.Int64), nil
}

func (d *withdrawDao) RecvNewWithdraw(withdraw entities.BaseWithdraw) (int64, error) {
	return insertTemplate((*baseDao)(unsafe.Pointer(d)), "RecvNewWithdraw", []interface {} {
		withdraw.Id,
		withdraw.Address,
		withdraw.Amount,
		withdraw.Asset,
	})
}

func (d *withdrawDao) WithdrawIntoStable(txHash string) (int64, error) {
	return updatePartsTemplate((*baseDao)(unsafe.Pointer(d)), "WithdrawIntoStable",
		[]interface {} { txHash }, nil)
}

func (d *withdrawDao) WithdrawIntoChain(txHash string, height uint64, txIndex int) (int64, error) {
	return updatePartsTemplate((*baseDao)(unsafe.Pointer(d)), "WithdrawIntoChain",
		[]interface {} { txHash }, map[string]interface {} {
			"height": height,
			"tx_index": txIndex,
		})
}

func (d *withdrawDao) SentForTxHash(txHash string, id int) (int64, error) {
	return updateTemplate((*baseDao)(unsafe.Pointer(d)), "SentForTxHash",
		[]interface {} { id }, []interface {} { txHash })
}

func (d *withdrawDao) GetWithdrawId(txHash string) (int, error) {
	var result []map[string]interface {}
	var err error
	bd := (*baseDao)(unsafe.Pointer(d))
	conds := []interface {} { txHash }
	if result, err = selectTemplate(bd, "GetWithdrawId", conds); err != nil {
		return -1, err
	}

	if len(result) != 1 {
		return -1, utils.LogMsgEx(utils.ERROR, "找不到交易：%s的提币ID", txHash)
	}
	return int(*result[0]["id"].(*int32)), nil
}

func (d *withdrawDao) GetWithdraws(conds map[string]interface {}) ([]entities.DatabaseWithdraw, error) {
	bd := (*baseDao)(unsafe.Pointer(d))
	var result []map[string]interface {}
	var err error
	if result, err = selectPartsTemplate(bd, "GetWithdraws", conds); err != nil {
		return nil, err
	}

	var ret []entities.DatabaseWithdraw
	for _, entity := range result {
		var withdraw entities.DatabaseWithdraw
		withdraw.Id = int(*entity["id"].(*int32))
		withdraw.TxHash = string(*entity["tx_hash"].(*sql.RawBytes))
		withdraw.To = string(*entity["address"].(*sql.RawBytes))
		withdraw.Address = withdraw.To
		withdraw.Amount, err = strconv.ParseFloat(string(*entity["amount"].(*sql.RawBytes)), 64)
		if err != nil {
			panic(utils.LogMsgEx(utils.ERROR, "解析交易金额失败：%v", err))
		}
		withdraw.Asset = string(*entity["asset"].(*sql.RawBytes))
		withdraw.Height = uint64(entity["height"].(*sql.NullInt64).Int64)
		withdraw.TxIndex = int(entity["tx_index"].(*sql.NullInt64).Int64)
		withdraw.CreateTime = entity["create_time"].(*mysql.NullTime).Time
		withdraw.UpdateTime = entity["update_time"].(*mysql.NullTime).Time
		ret = append(ret, withdraw)
	}
	return ret, nil
}

func (d *withdrawDao) CheckExistsById(id int) (bool, error) {
	bd := (*baseDao)(unsafe.Pointer(d))
	conds := []interface {} { id }
	var result []map[string]interface {}
	var err error
	if result, err = selectTemplate(bd, "CheckExistsById", conds); err != nil {
		return false, err
	}

	if len(result) != 1 {
		return false, utils.LogMsgEx(utils.ERROR, "数据库查询结果异常：COUNT返回无结果")
	}
	var ok bool
	var tmp interface {}
	if tmp, ok = result[0]["num"]; !ok {
		return false, utils.LogMsgEx(utils.ERROR, "COUNT结果没有指定键值：num")
	}
	return *tmp.(*int64) != 0, nil
}