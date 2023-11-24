// provides abstraction over reading a writing the lock file, as well as
// comparing it to the configured installed packages in the configuration
package lock

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/xnacly/mehr/types"
)

var path = LookUp()
var lock *types.LockFile = &types.LockFile{
	Packages: map[string]map[string]*types.Package{},
}

func init() {
	toml.DecodeFile(path, lock)
}

// path to the lock file
func LookUp() string {
	confHome, _ := os.UserConfigDir()
	return filepath.Join(confHome, "mehr", "lock.toml")
}

func All(lock *types.LockFile) map[string][]string {
	perm := map[string][]string{}
	// iterate over package managers
	for mgr, v := range lock.Packages {
		t := make([]string, 0, len(v))
		for pkgkey := range v {
			t = append(t, pkgkey)
		}
		perm[mgr] = t
	}
	return perm
}

// return all packages permanently installed on the system
func Permanent(config *types.Configuration, lock *types.LockFile) map[string][]string {
	perm := map[string][]string{}
	// iterate over package managers
	for mgr, v := range lock.Packages {
		t := make([]string, 0, len(v))
		for pkgkey := range v {
			if _, ok := config.Packages[mgr][pkgkey]; ok {
				t = append(t, pkgkey)
			}
		}
		perm[mgr] = t
	}
	return perm
}

func Len(packages map[string][]string) int {
	r := 0
	for _, l := range packages {
		if len(l) == 0 {
			continue
		}
		r += len(l)
	}
	return r
}

// return all temporary installed packages on the system
func Temporary(config *types.Configuration, lock *types.LockFile) map[string][]string {
	temp := map[string][]string{}
	// iterate over package managers
	for mgr, v := range lock.Packages {
		t := make([]string, 0, len(v))
		for pkgkey := range v {
			if _, ok := config.Packages[mgr][pkgkey]; !ok {
				t = append(t, pkgkey)
			}
		}
		temp[mgr] = t
	}
	return temp
}

// reads the lockfile from disk on first call, every call afterwards omits the
// disk interaction and reads the cached lock file
func Get() *types.LockFile {
	return lock
}

func writeToDisk() error {
	val, err := os.Stat(path)
	if err == nil && val.IsDir() {
		return fmt.Errorf("Lockfile %q is a directory, how did that happen?", path)
	} else {
		os.MkdirAll(filepath.Dir(path), 0777)
	}

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("Failed to open %q, %s", path, err)
	}
	defer file.Close()

	e := toml.NewEncoder(file)
	return e.Encode(lock)
}

func UpdateTimeStamp() error {
	lock.LastUpdate = time.Now()
	return writeToDisk()
}

// adds entry to the lock file, writes it to disk and updates the cached lock file
func AddPackage(name, packagemanager string, entry *types.Package) error {
	if _, ok := lock.Packages[packagemanager]; !ok {
		lock.Packages[packagemanager] = make(map[string]*types.Package)
	}
	lock.Packages[packagemanager][name] = entry
	return writeToDisk()
}

func RemovePackage(name, packagemanager string) error {
	delete(lock.Packages[packagemanager], name)
	return writeToDisk()
}
