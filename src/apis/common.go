package apis

import (
	"net/http"
	"fmt"
	"utils"
	"regexp"
	"strings"
	"reflect"
	"encoding/json"
)

const WithdrawPath = "/api/withdraw"
const WpLen = len(WithdrawPath)
const DepositPath = "/api/deposit"
const DpLen = len(DepositPath)

type respVO struct {
	Code int			`json:code`
	Msg string			`json:message`
	Data interface {}	`json:data`
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
		} else {
			utils.LogIdxEx(utils.WARNING, 36, handle.Method, req.Method)
			var resp respVO
			resp.Code = 404
			resp.Msg = fmt.Sprintf(utils.GetIdxMsg("W0036"), handle.Method, req.Method)
			respJSON , _:= json.Marshal(resp)
			w.Write(respJSON)
		}
	}
}

func RootHandler(w http.ResponseWriter, req *http.Request) {
	utils.LogMsgEx(utils.INFO, "%s\t\t%s", req.Method, req.RequestURI)
	switch {
	case len(req.RequestURI) >= WpLen && req.RequestURI[:WpLen] == WithdrawPath:
		subHandler(w, req, wdRouteMap)
	case len(req.RequestURI) >= DpLen && req.RequestURI[:DpLen] == DepositPath:
		subHandler(w, req, dpRouteMap)
	}
}