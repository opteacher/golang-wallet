package main

import (
	"fmt"
	"utils"
	"log"
)

type TestObs struct {
}

func (o *TestObs) BeforeTurn(s *utils.Status, tgtStt int) {
	log.Printf("Before turn: %d, to %d\n", s.Current(), tgtStt)
}

func (o *TestObs) AfterTurn(s *utils.Status, srcStt int) {
	log.Printf("After turn: %d, from %d\n", s.Current(), srcStt)
}

func main() {
	log.SetFlags(log.Lshortfile)

	//Test status and observer
	var err error
	fmt.Println()
	o := TestObs {}
	const (
		NONE = iota
		INIT
		START
		UNEXISTS
	)
	a := utils.Status {
		AllStatus:	[]int { NONE, INIT, START },
	}
	a.RegAsObs(&o)
	log.Println(a.Current())

	a.TurnTo(START)
	log.Println(a.Current())

	if _, err = a.TurnTo(UNEXISTS); err != nil {
		log.Println(err)
	}
}
