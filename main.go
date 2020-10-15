package main

import (
	"cornstone/aliases/e"
	"cornstone/cmd"
	"log"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatal(e.P(err))
	}
}
