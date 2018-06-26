package apis

import (
	"net/http"
)

func newAddress(w http.ResponseWriter, req *http.Request) []byte {
	// @_@: 币种名应该从path中获得
	//coinSet := utils.GetConfig().GetCoinSettings()
	//rpc := rpcs.GetRPC(coinSet.Name)
	return []byte("abcd")
}

var dpRouteMap = map[string]api {
	"^/api/deposit/([A-Z]{3,})/address":	{ newAddress, "GET" },
}