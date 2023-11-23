package types

import "time"

type LockFile struct {
	LastUpdate time.Time                      `toml:"last-update"`
	Packages   map[string]map[string]*Package `toml:"packages"`
}
