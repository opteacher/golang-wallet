package apis

import (
	"net/http"
	"fmt"
	"regexp"
	"reflect"
)

func test() {

}

func queryProcess(w http.ResponseWriter, req *http.Request) []byte {
	fmt.Println("abcd")
	return []byte("ttt")
}

var routeMap = map[string]interface {} {
	"^/api/withdraw/([A-Z]{3,})/process$": queryProcess,
}

func WithdrawHandler(w http.ResponseWriter, req *http.Request) {
	for route, handle := range routeMap {
		re := regexp.MustCompile(route)
		if re.MatchString(req.RequestURI) {
			a := reflect.ValueOf(handle).Call([]reflect.Value {
				reflect.ValueOf(w), reflect.ValueOf(req),
			})
			fmt.Println(a)
			w.Write(a[0].Bytes())
		}
	}
}