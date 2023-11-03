package main

import (
	"log"

	"github.com/xnacly/mehr/config"
)

func main() {
	confPath := config.LookUp()
	conf, err := config.Get(confPath)
	if err != nil {
		log.Printf("%s, falling back to default", err)
	}
	log.Println(conf.PackageManager)
}
