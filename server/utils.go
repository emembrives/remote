package main

import (
	"log"
)

func logFatalOnError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
