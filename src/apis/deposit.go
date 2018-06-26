package apis

import (
	"net/http"
	"regexp"
	"rpcs"
	"encoding/json"
)

const newAddressPath = "^/api/deposit/([A-Z]{3,})/address$"

func newAddress(w http.ResponseWriter, req *http.Request) []byte {
	re := regexp.MustCompile(newAddressPath)
	params := re.FindStringSubmatch(req.RequestURI)[1:]
	coinName := params[0]
	rpc := rpcs.GetRPC(coinName)

	var address string
	var err error
	var resp respVO
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

var dpRouteMap = map[string]api {
	newAddressPath:	{ newAddress, "GET" },
}