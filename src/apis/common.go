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
	"rpcs"
	"io/ioutil"
	"net"
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

func subHandler(w http.ResponseWriter, req *http.Request, routeMap map[string]interface {}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	uri := strings.Split(req.RequestURI, "?")
	if len(uri) == 2 {
		req.RequestURI = uri[0]
		req.ParseForm()
	}
	for route, handle := range routeMap {
		reqGrp := strings.Split(route, " ")
		if len(reqGrp) != 2 {
			continue
		}
		method := reqGrp[0]
		path := reqGrp[1]
		re := regexp.MustCompile(path)
		if !re.MatchString(req.RequestURI) { continue }
		if strings.ToUpper(req.Method) == method {
			a := reflect.ValueOf(handle).Call([]reflect.Value {
				reflect.ValueOf(w), reflect.ValueOf(req),
			})
			w.Write(a[0].Bytes())
			return
		} else {
			utils.LogIdxEx(utils.WARNING, 36, method, req.Method)
			//var resp RespVO
			//resp.Code = 405
			//resp.Msg = fmt.Sprintf(utils.GetIdxMsg("W0036"), method, req.Method)
			//respJSON , _:= json.Marshal(resp)
			//w.Write(respJSON)
			//return
		}
	}
	utils.LogIdxEx(utils.WARNING, 37, req.RequestURI)
	var resp RespVO
	resp.Code = 404
	resp.Msg = fmt.Sprintf(utils.GetIdxMsg("W0037"), req.RequestURI)
	respJSON , _:= json.Marshal(resp)
	w.Write(respJSON)
}

func HttpHandler(w http.ResponseWriter, req *http.Request) {
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

func SocketHandler(conn net.Conn) {
	buffer := make([]byte, 2048)
	for {
		var n int
		var err error
		if n, err = conn.Read(buffer); err != nil {
			utils.LogMsgEx(utils.ERROR, "SOCKET连接错误：%v", err)
			return
		}
		utils.LogMsgEx(utils.INFO, "SOCKET\t%s", string(buffer[:n]))
	}
}

const getProcessPath	= "^/api/common/([A-Z]{3,})/process/([a-zA-Z0-9]{1,})"
const transferPath		= "^/api/common/([A-Z]{3,})/transfer$"

var cmRouteMap = map[string]interface {} {
	fmt.Sprintf("%s %s", http.MethodGet, getProcessPath): queryProcess,
	fmt.Sprintf("%s %s", http.MethodPost, transferPath): transfer,
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

type transactionReq struct {
	From string		`json:"from"`
	To string		`json:"to"`
	Amount float64	`json:"amount"`
}

func transfer(w http.ResponseWriter, req *http.Request) []byte {
	var resp RespVO
	re := regexp.MustCompile(transferPath)
	params := re.FindStringSubmatch(req.RequestURI)[1:]
	if len(params) == 0 {
		resp.Code = 500
		resp.Msg = "需要指定币种的名字"
		ret, _ := json.Marshal(resp)
		return ret
	}

	// 参数解析
	var body []byte
	var err error
	if body, err = ioutil.ReadAll(req.Body); err != nil {
		utils.LogMsgEx(utils.WARNING, "解析请求体错误：%v", err)
		resp.Code = 500
		resp.Msg = err.Error()
		ret, _ := json.Marshal(resp)
		return ret
	}
	defer req.Body.Close()

	utils.LogMsgEx(utils.INFO, "收到交易请求：%s", string(body))

	var txReq transactionReq
	if err = json.Unmarshal(body, &txReq); err != nil {
		utils.LogIdxEx(utils.WARNING, 38, err)
		resp.Code = 500
		resp.Msg = err.Error()
		ret, _ := json.Marshal(resp)
		return ret
	}

	rpc := rpcs.GetRPC(params[0])
	var txHash string
	tradePwd := utils.GetConfig().GetCoinSettings().TradePassword
	if txHash, err = rpc.SendTransaction(txReq.From, txReq.To, txReq.Amount, tradePwd); err != nil {
		utils.LogMsgEx(utils.ERROR, "发送交易失败：%v", err)
		resp.Code = 500
		resp.Msg = err.Error()
		ret, _ := json.Marshal(resp)
		return ret
	}

	resp.Code = 200
	resp.Data = txHash
	ret, _ := json.Marshal(resp)
	return []byte(ret)
}