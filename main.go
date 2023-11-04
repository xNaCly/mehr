package main

import (
	"fmt"
	"log"

	"github.com/xnacly/mehr/config"
	l "github.com/xnacly/mehr/log"
	packagemanagers "github.com/xnacly/mehr/package_managers"
)

func main() {
	log.SetFlags(log.Ltime)
	confPath := config.LookUp()
	_, err := config.Get(confPath)
	if err != nil {
		l.Error("%s, falling back to default", err)
	}
	if manager, ok := packagemanagers.Get(); ok {
		l.Info("detected package manager %q", manager.Name)
		if err := manager.Update(); err != nil {
			l.Error("%q, printing stdout and stderr: ", err)
			fmt.Println(manager.Buffer)
		}
	}
}
