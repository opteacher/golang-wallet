package apis

import (
	"net/http"
	"fmt"
)

var wdRouteMap = map[string]api {
	"^/api/withdraw/([A-Z]{3,})/process$":	{ queryProcess, "GET" },
	"^/api/withdraw/([A-Z]{3,})":			{ doWithdraw, "POST" },
}

func queryProcess(w http.ResponseWriter, req *http.Request) []byte {
	fmt.Println("abcd")
	return []byte("ttt")
}

func doWithdraw(w http.ResponseWriter, req *http.Request) []byte {
	return []byte("abcd")
}