package main

import (
	"encoding/json"
	"log"
	"bytes"
	"net/http"
	"io/ioutil"
)

const URL = "http://18.144.17.127:8545"
type ReqBody struct {
	Method string	`json:method`
	Params []string	`json:params`
	Id string		`json:id`
}

func main() {
	log.SetFlags(log.Lshortfile)

	// Request from blockchain
	reqBody := ReqBody { "eth_blockNumber", []string {}, "latest" }
	reqStr, err := json.Marshal(reqBody)
	if err != nil {
		log.Fatal(err)
		return
	}
	reqBuf := bytes.NewBuffer([]byte(reqStr))
	res, err := http.Post(URL, "application/json", reqBuf)
	defer res.Body.Close()

	// Parse response body
	bodyStr, err := ioutil.ReadAll(res.Body)
	log.Println(string(bodyStr))
}
