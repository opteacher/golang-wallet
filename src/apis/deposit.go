package apis

import (
	"net/http"
	"regexp"
	"rpcs"
	"encoding/json"
)

const newAddressPath	= "^/api/deposit/([A-Z]{3,})/address$"
const getHeightPath		= "^/api/deposit/([A-Z]{3,})/height$"

var dpRouteMap = map[string]api {
	newAddressPath:	{ newAddress, "GET" },
	getHeightPath:	{ queryHeight, "GET"},
}

func newAddress(w http.ResponseWriter, req *http.Request) []byte {
	var resp RespVO
	re := regexp.MustCompile(newAddressPath)
	params := re.FindStringSubmatch(req.RequestURI)[1:]
	if len(params) == 0 {
		resp.Code = 500
		resp.Msg = "需要指定币种的名字"
		ret, _ := json.Marshal(resp)
		return ret
	}

	coinName := params[0]
	rpc := rpcs.GetRPC(coinName)

	var address string
	var err error
	if address, err = rpc.GetNewAddress(); err != nil {
		resp.Code = 500
		resp.Msg = err.Error()
		ret, _ := json.Marshal(resp)
		return ret
	}

	resp.Code = 200
	resp.Data = address
	ret, _ := json.Marshal(resp)
	return ret
}

func queryHeight(w http.ResponseWriter, req *http.Request) []byte {
	var resp RespVO
	re := regexp.MustCompile(getHeightPath)
	params := re.FindStringSubmatch(req.RequestURI)[1:]
	if len(params) == 0 {
		resp.Code = 500
		resp.Msg = "需要指定币种的名字"
		ret, _ := json.Marshal(resp)
		return ret
	}

	coinName := params[0]
	rpc := rpcs.GetRPC(coinName)

	var height uint64
	var err error
	if height, err = rpc.GetCurrentHeight(); err != nil {
		resp.Code = 500
		resp.Msg = err.Error()
		ret, _ := json.Marshal(resp)
		return ret
	}

	resp.Code = 200
	resp.Data = height
	ret, _ := json.Marshal(resp)
	return ret
}