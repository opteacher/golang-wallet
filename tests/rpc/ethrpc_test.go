package rpc

import (
	"log"
	"rpcs"
	"testing"
	"fmt"
)

func TestGetTransactions(t *testing.T) {
	log.SetFlags(log.Lshortfile)

	//Test RPC
	fmt.Println(rpcs.GetRPC("ETH").GetTransactions(0))
}

func TestSendTransaction(t *testing.T) {
	var err error
	var txHash string
	if txHash, err = rpcs.GetRPC("ETH").SendFrom(
		"0x47e8e8e49a8c1c308e84439f2c55ef18710f5ed6",
		10.23456); err != nil {
		log.Fatal(err)
	}
	log.Println(txHash)
}

func TestNewAddress(t *testing.T) {
	var addr string
	var err error
	if addr, err = rpcs.GetRPC("ETH").GetNewAddress(); err != nil {
		log.Fatal(err)
	} else {
		log.Println(addr)
	}
}

func TestGetTxExistsHeight(t *testing.T) {
	fmt.Println(rpcs.GetRPC("ETH").GetTxExistsHeight("abcd"))
}