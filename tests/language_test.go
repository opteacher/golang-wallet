package tests

import (
	"log"
	"fmt"
	"reflect"
	"testing"
)

type itfc interface {
	Test() string
}

type test struct {

}

func (t *test) Test() string {
	return "abcd"
}

func TestLang(t *testing.T) {
	log.SetFlags(log.Lshortfile)
	fmt.Println("abcd")
	var c float64 = 58500000000000000000
	log.Println(reflect.TypeOf(c).Name())

	var a itfc
	var b = new(test)
	a = b
	fmt.Println(a.Test())
}
