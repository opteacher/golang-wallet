package dao

import (
	"sync"
	"entities"
	"database/sql"
	"utils"
	"unsafe"
)

type processDao struct {
	baseDao
	sync.Once
}

var _processDao *processDao

func GetProcessDAO() *processDao {
	if _processDao == nil {
		_processDao = new(processDao)
		_processDao.Once = sync.Once {}
		_processDao.Once.Do(func() {
			_processDao.create("process")
		})
	}
	return _processDao
}

func (d *processDao) SaveProcess(process *entities.DatabaseProcess) (int64, error) {
	var keys = []string { "tx_hash" }
	var vals = []interface {} { process.TxHash }
	if process.Asset != "" {
		keys = append(keys, "asset")
		vals = append(vals, process.Asset)
	}
	if process.Type != "" {
		keys = append(keys, "`type`")
		vals = append(vals, process.Type)
	}
	if process.Height != 0 {
		keys = append(keys, "height")
		vals = append(vals, process.Height)
	}
	if process.CompleteHeight != 0 {
		keys = append(keys, "complete_height")
		vals = append(vals, process.CompleteHeight)
	}
	if utils.StrArrayContains(entities.Processes, process.Process) {
		keys = append(keys, "process")
		vals = append(vals, process.Process)
	}
	if !process.Cancelable {
		keys = append(keys, "cancelable")
		vals = append(vals, 0)
	}
	return saveTemplate((*baseDao)(unsafe.Pointer(d)),
		"CheckProcsExists", "AddProcess", "UpdateProcessByHash",
			[]interface {} { process.TxHash }, vals, keys)
}

func (d *processDao) QueryProcess(asset string, txHash string) (entities.DatabaseProcess, error) {
	var ret entities.DatabaseProcess
	var result []map[string]interface {}
	var err error
	conds := []interface {} { asset, txHash }
	bd := (*baseDao)(unsafe.Pointer(d))
	if result, err = selectTemplate(bd, "QueryProcess", conds); err != nil {
		return ret, err
	}
	if len(result) != 1 {
		return ret, utils.LogMsgEx(utils.ERROR, "无法找到指定的进度任务：%s", txHash)
	}

	entity := result[0]
	ret.TxHash = txHash
	ret.Asset = asset
	ret.Type = string(*entity["type"].(*sql.RawBytes))
	ret.Height = uint64(entity["height"].(*sql.NullInt64).Int64)
	ret.CurrentHeight = uint64(entity["current_height"].(*sql.NullInt64).Int64)
	ret.CompleteHeight = uint64(entity["complete_height"].(*sql.NullInt64).Int64)
	ret.Process = string(*entity["process"].(*sql.RawBytes))
	ret.Cancelable = *entity["cancelable"].(*int8) != 0
	return ret, nil
}