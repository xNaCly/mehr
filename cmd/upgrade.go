package cmd

import (
	"time"

	"github.com/spf13/cobra"
	"github.com/xnacly/mehr/config"
	"github.com/xnacly/mehr/lock"
	l "github.com/xnacly/mehr/log"
	"github.com/xnacly/mehr/pkgmgr"
	"github.com/xnacly/mehr/types"
)

func init() {
	rootCmd.AddCommand(upgradeCmd)
}

var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade packages",
	Example: `
    upgrade neovim
    upgrade neovim kitty
    upgrade`,
	Long: `Upgrade a single, multiple or all currently installed packages via: 

Upgrade a single package:

    upgrade [package]

Upgrade multiple packages:

    upgrade [package...]

Upgrade all packages in the current configuration:

    upgrade

Errors on attempting to upgrade packages not found in the lock file.
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

		// only update every hour
		if time.Since(lock.Get().LastUpdate).Hours() >= 1 {
			err := lock.UpdateTimeStamp()
			if err != nil {
				l.Errorf("Failed to update the lock file update timestamp")
			}
			err = manager.Update()
			if err != nil {
				l.Errorf("Failed to upgrade package manager repositories: %w", err)
			}
		} else {
			l.Infof("Last sync less than an hour ago, skipping repo updates")
		}

		if len(args) == 0 {
			err, amount := manager.Upgrade(lock.Get().Packages)
			if err != nil {
				l.Error("failed to upgrade packages", err)
			} else if amount > 0 {
				l.Infof("Upgraded %d packages", len(conf.Packages))
			} else {
				l.Infof("Did nothing, exiting")
			}
		} else {
			pkgs := map[string]*types.Package{}
			for _, pkg := range args {
				pkgs[pkg] = &types.Package{}
			}
			err, amount := manager.Install(pkgs)
			if err != nil {
				l.Error("failed to upgrade packages", err)
			} else if amount > 0 {
				l.Infof("Upgrade %d packages", len(conf.Packages))
			}
		}
	},
}
