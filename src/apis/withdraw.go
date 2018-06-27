package apis

import (
	"net/http"
	"regexp"
	"encoding/json"
)

var wdRouteMap = map[string]api {
	"^/api/withdraw/([A-Z]{3,})":	{ doWithdraw, "POST" },
}

func doWithdraw(w http.ResponseWriter, req *http.Request) []byte {
	var resp RespVO
	re := regexp.MustCompile(getHeightPath)
	params := re.FindStringSubmatch(req.RequestURI)[1:]
	if len(params) == 0 {
		resp.Code = 500
		resp.Msg = "需要指定币种的名字"
		ret, _ := json.Marshal(resp)
		return ret
	}

	return []byte("abcd")
}