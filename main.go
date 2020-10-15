package main

import (
	"cornstone/aliases/e"
	"cornstone/cmd"
	"log"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatalln(e.P(err))
	}
}
