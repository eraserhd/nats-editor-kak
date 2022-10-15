package main

import (
	"log"
	"os"

	"github.com/plugbench/kakoune-pluggo/service"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Syntax is `kakoune-pluggo-daemon SESSION`.")
	}
	session := os.Args[1]

	es, err := service.New(session)
	if err != nil {
		log.Fatal(err)
	}
	if err := es.Run(); err != nil {
		log.Fatal(err)
	}
	return
}
