package cmd

import (
	"github.com/xnacly/mehr/config"

	l "github.com/xnacly/mehr/log"
)

func Root() {
	confPath := config.LookUp()

	conf, err := config.Get(confPath)

	if err != nil {
		l.Error("%s", err)
	}

	l.Info("Got %q for Configuration.PackageManager", conf.PackageManager)
}
