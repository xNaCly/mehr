package main

import (
	"log"

	"github.com/xnacly/mehr/cmd"
)

func main() {
	log.SetFlags(log.Ltime)
	cmd.Root()
}
