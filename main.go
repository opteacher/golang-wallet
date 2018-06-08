package main

import (
	"encoding/json"
	"net/http"
	"fmt"
	"strings"
	"io/ioutil"
	"managers"
)

const URL = "http://18.144.17.127:8545"
const Method = "POST"

type ReqBody struct {
	Method string	`json:"method"`
	Params []string	`json:"params"`
	Id int			`json:"id"`
}

func main() {
	managers.Test()

	// 包装请求体
	reqBody := &ReqBody { "eth_blockNumber", []string {}, 1 }
	reqBodyBtAry, err := json.Marshal(reqBody)
	if err != nil {
		fmt.Println("error: ", err)
		return
	} else {
		fmt.Println(string(reqBodyBtAry))
	}
	reqBodyStr := string(reqBodyBtAry)

	// 发送请求
	resp, err := http.Post(URL, "application/json", strings.NewReader(reqBodyStr))
	if err != nil {
		fmt.Println("error: ", err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("error: ", err)
		return
	}
	fmt.Println(string(body))
}