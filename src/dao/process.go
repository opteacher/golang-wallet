package dao

import (
	"sync"
	"entities"
	"database/sql"
	"databases"
	"utils"
	"fmt"
	"strings"
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

func (dao *processDao) SaveProcess(process *entities.DatabaseProcess) (int64, error) {
	var db *sql.DB
	var err error
	if db, err = databases.ConnectMySQL(); err != nil {
		panic(utils.LogIdxEx(utils.ERROR, 0010, err))
	}

	var checkSQL string
	var ok bool
	if checkSQL, ok = dao.sqls["CheckProcsExists"]; !ok {
		return 0, utils.LogIdxEx(utils.ERROR, 0011, "CheckProcsExists")
	}

	var rows *sql.Rows
	if rows, err = db.Query(checkSQL, process.TxHash); err != nil {
		panic(utils.LogIdxEx(utils.ERROR, 0013, err))
	}
	defer rows.Close()

	var num int
	if rows.Next() {
		if err = rows.Scan(&num); err != nil {
			return 0, utils.LogIdxEx(utils.ERROR, 0014, err)
		}
	}

	var result sql.Result
	if num == 0 {
		var insertSQL string
		if insertSQL, ok = dao.sqls["AddProcess"]; !ok {
			return 0, utils.LogIdxEx(utils.ERROR, 0011, "AddProcess")
		}

		if result, err = db.Exec(insertSQL, []interface {} {
			process.TxHash,
			process.Type,
			process.Height,
			process.CompleteHeight,
			process.Process,
			process.Cancelable,
		}...); err != nil {
			panic(utils.LogIdxEx(utils.ERROR, 0012, err))
		}
		return result.RowsAffected()
	} else {
		var updateSQL string
		if updateSQL, ok = dao.sqls["UpdateProcessByHash"]; !ok {
			return 0, utils.LogIdxEx(utils.ERROR, 0011, "UpdateProcessByHash")
		}

		var sets []string
		var vals []interface {}
		if process.Height != 0 {
			sets = append(sets, "height=?")
			vals = append(vals, process.Height)
		}
		if process.CompleteHeight != 0 {
			sets = append(sets, "complete_height=?")
			vals = append(vals, process.CompleteHeight)
		}
		if utils.StrArrayContains(entities.Processes, process.Process) {
			sets = append(sets, "process=?")
			vals = append(vals, process.Process)
		}
		if !process.Cancelable {
			sets = append(sets, "cancelable=?")
			vals = append(vals, 0)
		}

		if len(sets) == 0 {
			return 0, utils.LogIdxEx(utils.ERROR, 0030, nil)
		}

		vals = append(vals, process.TxHash)
		updateSQL = fmt.Sprintf(updateSQL, strings.Join(sets, ","))
		if result, err = db.Exec(updateSQL, vals...); err != nil {
			panic(utils.LogIdxEx(utils.ERROR, 0012, err))
		}
		return result.RowsAffected()
	}
}