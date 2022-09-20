package main

import (
	"log"

	"github.com/eraserhd/nats-editor-kak/service"
)

func main() {
        es, err := service.New()
        if err != nil {
                log.Fatal(err)
        }
        if err := es.Run(); err != nil {
                log.Fatal(err)
        }
}
