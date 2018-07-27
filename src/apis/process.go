package apis

import (
	"strconv"
	"strings"
	"entities"
	"encoding/json"
	"dao"
	"net/http"
	"regexp"
	"fmt"
)

const getProcessPath	= "^/api/process/([A-Z]{3,})/txid/([a-zA-Z0-9]{1,})"

var pcsRouteMap = map[string]interface {} {
	fmt.Sprintf("%s %s", http.MethodGet, getProcessPath): queryProcess,
}

func queryProcess(w http.ResponseWriter, req *http.Request) []byte {
	var resp RespVO
	re := regexp.MustCompile(getProcessPath)
	params := re.FindStringSubmatch(req.RequestURI)[1:]
	if len(params) == 0 {
		resp.Code = 500
		resp.Msg = "需要指定币种的名字"
		ret, _ := json.Marshal(resp)
		return ret
	}
	if len(params) == 1 {
		resp.Code = 500
		resp.Msg = "需要指定查询的操作ID或交易哈希"
		ret, _ := json.Marshal(resp)
		return ret
	}

	coinName := params[0]
	txId := params[1]

	var err error
	var id int64 = -1
	var typ string
	if id, err = strconv.ParseInt(txId, 10, 64); err == nil {
		typ = strings.ToUpper(req.URL.Query().Get("type"))
		if typ != entities.WITHDRAW && typ != entities.DEPOSIT {
			resp.Code = 500
			resp.Msg = "用操作id查询进度，需附带操作类型：WITHDRAW/DEPOSIT"
			ret, _ := json.Marshal(resp)
			return ret
		}
	}

	var process entities.DatabaseProcess
	if id == -1 {
		if process, err = dao.GetProcessDAO().QueryProcessByTypAndId(coinName, typ, int(id)); err != nil {
			resp.Code = 500
			resp.Msg = err.Error()
			ret, _ := json.Marshal(resp)
			return ret
		}
	} else {
		if process, err = dao.GetProcessDAO().QueryProcessByTxHash(coinName, txId); err != nil {
			resp.Code = 500
			resp.Msg = err.Error()
			ret, _ := json.Marshal(resp)
			return ret
		}
	}

	resp.Code = 200
	resp.Data = process
	ret, _ := json.Marshal(resp)
	return ret
}