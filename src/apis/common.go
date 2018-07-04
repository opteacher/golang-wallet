package apis

import (
	"net/http"
	"regexp"
	"encoding/json"
	"entities"
	"dao"
	"fmt"
	"utils"
	"strings"
	"reflect"
	"strconv"
)

const WithdrawPath = "/api/withdraw"
const WpLen = len(WithdrawPath)
const DepositPath = "/api/deposit"
const DpLen = len(DepositPath)
const CommonPath = "/api/common"
const CmLen = len(CommonPath)

type RespVO struct {
	Code int			`json:"code"`
	Msg string			`json:"message"`
	Data interface {}	`json:"data"`
}

type api struct {
	Func interface {}
	Method string
}

func subHandler(w http.ResponseWriter, req *http.Request, routeMap map[string]api) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	for route, handle := range routeMap {
		re := regexp.MustCompile(route)
		if !re.MatchString(req.RequestURI) { continue }
		if strings.ToUpper(req.Method) == handle.Method {
			a := reflect.ValueOf(handle.Func).Call([]reflect.Value {
				reflect.ValueOf(w), reflect.ValueOf(req),
			})
			w.Write(a[0].Bytes())
			return
		} else {
			utils.LogIdxEx(utils.WARNING, 36, handle.Method, req.Method)
			var resp RespVO
			resp.Code = 405
			resp.Msg = fmt.Sprintf(utils.GetIdxMsg("W0036"), handle.Method, req.Method)
			respJSON , _:= json.Marshal(resp)
			w.Write(respJSON)
			return
		}
	}
	utils.LogIdxEx(utils.WARNING, 37, req.RequestURI)
	var resp RespVO
	resp.Code = 404
	resp.Msg = fmt.Sprintf(utils.GetIdxMsg("W0037"), req.RequestURI)
	respJSON , _:= json.Marshal(resp)
	w.Write(respJSON)
}

func RootHandler(w http.ResponseWriter, req *http.Request) {
	utils.LogMsgEx(utils.INFO, "%s\t\t%s", req.Method, req.RequestURI)
	switch {
	case len(req.RequestURI) >= WpLen && req.RequestURI[:WpLen] == WithdrawPath:
		subHandler(w, req, wdRouteMap)
	case len(req.RequestURI) >= DpLen && req.RequestURI[:DpLen] == DepositPath:
		subHandler(w, req, dpRouteMap)
	case len(req.RequestURI) >= CmLen && req.RequestURI[:CmLen] == CommonPath:
		subHandler(w, req, cmRouteMap)
	default:
		utils.LogIdxEx(utils.WARNING, 37, req.RequestURI)
		var resp RespVO
		resp.Code = 404
		resp.Msg = fmt.Sprintf(utils.GetIdxMsg("W0037"), req.RequestURI)
		respJSON , _:= json.Marshal(resp)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.Write(respJSON)
	}
}

const getProcessPath	= "^/api/common/([A-Z]{3,})/process/([a-zA-Z0-9]{1,})"

var cmRouteMap = map[string]api{
	getProcessPath: {queryProcess, http.MethodGet },
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