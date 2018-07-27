package tests

import (
	"log"
	"fmt"
	"reflect"
	"testing"
	"math/big"
	"math"
	"strings"
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

	reflect.ArrayOf(10, reflect.TypeOf(test {}))
	fmt.Println()

	abcd := []int {1}
	abcd = append(abcd[:0], abcd[1:]...)
	fmt.Println(abcd)

	var ttttt = big.NewFloat(0)
	var y = 0
	var err error
	if ttttt, y, err = ttttt.Parse("6194049F30F7200000", 16); err != nil {
		log.Fatal(err)
	} else {
		log.Println(ttttt.Mul(ttttt, big.NewFloat(math.Pow10(-18))).String())
		log.Println(y)
	}

	aa := 1
	bb := 2
	cc := 3
	switch {
	case aa == 1:
		fmt.Println("aa")
	fallthrough
	case bb == 2:
		fmt.Println("bb")
	fallthrough
	case cc == 3:
		fmt.Println("cc")
	default:
		fmt.Println("dd")
	}

	fmt.Println(strings.Split("abcd", "."))
}
