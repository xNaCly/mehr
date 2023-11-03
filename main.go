package main

import (
	"log"

	packagemanagers "github.com/xnacly/mehr/package_managers"
)

func main() {
	if pkgMgr, ok := packagemanagers.Get(); ok {
		log.Println(pkgMgr.Name)
	}
}
