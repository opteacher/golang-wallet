package rpc

import (
	"testing"
	"fmt"
	"rpcs"
)

func TestBtcGetCurrentHeight(t *testing.T) {
	fmt.Println(rpcs.GetRPC("BTC").GetCurrentHeight())
}

func TestBtcGetTransactions(t *testing.T) {
	fmt.Println(rpcs.GetRPC("BTC").GetTransactions(1))
}

func TestBtcGetTxExistsHeight(t *testing.T) {
	fmt.Println(rpcs.GetRPC("BTC").GetTxExistsHeight(
		"3334ed8aa4a1c215f7c3ce6445dad14b2cd66cd004a1df39c9d6b6aad68d8edd"))
}

func TestGetNewAddress(t *testing.T) {
	fmt.Println(rpcs.GetRPC("BTC").GetNewAddress())
}

func TestBtcGetBalance(t *testing.T) {
	fmt.Println(rpcs.GetRPC("BTC").GetBalance(""))
}