package apis

import (
	"net/http"
	"regexp"
	"encoding/json"
	"fmt"
	"utils"
	"strings"
	"reflect"
	"net"
)

const WithdrawPath = "/api/withdraw"
const WpLen = len(WithdrawPath)
const DepositPath = "/api/deposit"
const DpLen = len(DepositPath)
const ProcessPath = "/api/process"
const PcsLen = len(ProcessPath)
const TestPath = "/api/test"
const TstLen = len(TestPath)

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
	case len(req.RequestURI) >= TstLen && req.RequestURI[:TstLen] == TestPath:
		subHandler(w, req, tstRouteMap)
	case len(req.RequestURI) >= PcsLen && req.RequestURI[:PcsLen] == ProcessPath:
		subHandler(w, req, pcsRouteMap)
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