package apis

import (
	"net/http"
	"fmt"
)

func WithdrawHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Println(req.RequestURI)
}