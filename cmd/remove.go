package cmd

import (
	"github.com/spf13/cobra"
	"github.com/xnacly/mehr/config"
	"github.com/xnacly/mehr/lock"
	l "github.com/xnacly/mehr/log"
	"github.com/xnacly/mehr/pkgmgr"
)

func init() {
	rootCmd.AddCommand(removeCmd)
}

var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove packages",
	Example: `
    remove neovim
    remove neovim kitty
    remove`,
	Long: `Remove a single, multiple or all currently installed packages via: 

Remove a single package:

    remove [package]

Remove multiple packages:

    remove [package...]

Remove all packages in the current configuration:

    remove

Errors on attempting to remove packages not found in the lock file
`,
	Run: func(cmd *cobra.Command, args []string) {
		configPath := config.LookUp()

		conf, err := config.Get(configPath)
		if err != nil {
			l.Errorf("Failed to get config: %q", err)
			return
		}

		var manager *pkgmgr.PackageManager
		if conf.PackageManager == "auto" || conf.PackageManager == "" {
			mgr, ok := pkgmgr.Get()

			if !ok {
				l.Error("Failed to find a package manager")
				return
			}

			manager = mgr
		} else {
			var err error
			mgr, err := pkgmgr.GetByName(conf.PackageManager)
			if err != nil {
				l.Error(err)
				return
			}
			manager = mgr
		}

		if len(args) == 0 {
			pkgs := make([]string, 0)
			for k, _ := range lock.Temporary(conf, lock.Get()) {
				pkgs = append(pkgs, k)
			}
			err, amount := manager.Remove(pkgs...)
			if err != nil {
				l.Error("failed to upgrade packages", err)
			} else if amount > 0 {
				l.Infof("Upgraded %d packages", len(conf.Packages))
			} else {
				l.Infof("Did nothing, exiting")
			}
		} else {
			err, amount := manager.Remove(args...)
			if err != nil {
				l.Error("failed to upgrade packages", err)
			} else if amount > 0 {
				l.Infof("Upgrade %d packages", len(conf.Packages))
			}
		}
	},
}