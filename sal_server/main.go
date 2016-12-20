package main

import (
	"log"

	"github.com/thisisfineio/sal"
)

func init() {}

func main() {
	err := sal.Run()
	if err != nil {
		log.Fatal(err)
	}
}
