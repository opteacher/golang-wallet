package apis

import (
	"strconv"
	"encoding/json"
	"dao"
	"net/http"
	"regexp"
	"fmt"
)

const processByTxidPath	= "^/api/process/([A-Z]{3,})/txid/([a-zA-Z0-9]{1,})"
const processByOpidPath	= "^/api/process/([A-Z]{3,})/type/(WITHDRAW|DEPOSIT)/id/([0-9]{1,})"

var pcsRouteMap = map[string]interface {} {
	fmt.Sprintf("%s %s", http.MethodGet, processByTxidPath): queryProcessByTxid,
	fmt.Sprintf("%s %s", http.MethodGet, processByOpidPath): queryProcessByOpid,
}

func queryProcessByTxid(w http.ResponseWriter, req *http.Request) []byte {
	var resp RespVO
	re := regexp.MustCompile(processByTxidPath)
	params := re.FindStringSubmatch(req.RequestURI)[1:]
	if len(params) == 0 {
		resp.Code = 500
		resp.Msg = "需要指定币种的名字"
		ret, _ := json.Marshal(resp)
		return ret
	}
	asset := params[0]
	if len(params) == 1 {
		resp.Code = 500
		resp.Msg = "需要指定查询的交易哈希"
		ret, _ := json.Marshal(resp)
		return ret
	}
	txId := params[1]

	if process, err := dao.GetProcessDAO().QueryProcessByTxHash(asset, txId); err != nil {
		resp.Code = 500
		resp.Msg = fmt.Sprintf("未找到指定交易：%v", err)
		ret, _ := json.Marshal(resp)
		return ret
	} else {
		resp.Code = 200
		resp.Data = process
		ret, _ := json.Marshal(resp)
		return ret
	}
}

func queryProcessByOpid(w http.ResponseWriter, req *http.Request) []byte {
	var resp RespVO
	var err error
	re := regexp.MustCompile(processByOpidPath)
	params := re.FindStringSubmatch(req.RequestURI)[1:]
	if len(params) == 0 {
		resp.Code = 500
		resp.Msg = "需要指定币种的名字"
		ret, _ := json.Marshal(resp)
		return ret
	}
	asset := params[0]
	if len(params) <= 2{
		resp.Code = 500
		resp.Msg = "需要指定查询的操作类型和操作id"
		ret, _ := json.Marshal(resp)
		return ret
	}
	typ := params[1]
	var id int
	if id, err = strconv.Atoi(params[2]); err != nil {
		resp.Code = 500
		resp.Msg = fmt.Sprintf("操作id必须是数字：%v", err)
		ret, _ := json.Marshal(resp)
		return ret
	}

	if process, err := dao.GetProcessDAO().QueryProcessByTypAndId(asset, typ, id); err != nil {
		resp.Code = 500
		resp.Msg = fmt.Sprintf("未找到指定交易：%v", err)
		ret, _ := json.Marshal(resp)
		return ret
	} else {
		resp.Code = 200
		resp.Data = process
		ret, _ := json.Marshal(resp)
		return ret
	}
}