package main

import (
	"log"
	"fmt"
	"reflect"
)

type itfc interface {
	Test() string
}

type test struct {

}

func (t *test) Test() string {
	return "abcd"
}

func main() {
	log.SetFlags(log.Lshortfile)
	fmt.Println("abcd")
	var t float64 = 58500000000000000000
	log.Println(reflect.TypeOf(t).Name())

	var a itfc
	var b = new(test)
	a = b
	fmt.Println(a.Test())

}
