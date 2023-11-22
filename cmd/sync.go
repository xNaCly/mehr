package cmd

import (
	"github.com/spf13/cobra"
	"github.com/xnacly/mehr/config"
	"github.com/xnacly/mehr/lock"
	l "github.com/xnacly/mehr/log"
	"github.com/xnacly/mehr/pkgmgr"
)

func init() {
	rootCmd.AddCommand(syncCmd)
}

var syncCmd = &cobra.Command{
	Use:     "sync",
	Short:   "Synchronise system state to configuration file",
	Example: `sync`,
	Long: `
Either forwards or resets the systems state to the state of the
configuration file, thus syncing both and removing the current lock file.

Packages not installed but found in the mehr configuration will be installed
upon running sync. If packages are installed on the system but aren't
reflected in the configuration (referred to as temporary packages), the system
can be synced to the configuration via mehr sync.
    `,
	Run: func(cmd *cobra.Command, args []string) {
		l.Infof("Syncing system state to match configuration file")
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

		tempPkg := lock.Temporary(conf, lock.Get())
		printPackages(tempPkg)
		if len(tempPkg) > 0 {
			l.Info("Removing temporally installed packages")
			if force, err := cmd.Flags().GetBool("force"); err != nil || !force {
				l.Errorf("Would permanently remove %d temporary packages, rerun with --force to continue", len(tempPkg))
				return
			} else {
				l.Warn("Got --force, mehr will remove packages and or configuration to match the configured system state")
			}

			packages := make([]string, 0)

			for k := range tempPkg {
				packages = append(packages, k)
			}

			err, amount := manager.Remove(packages...)
			if err != nil {
				l.Errorf("Failed to remove temporary packages: %w", err)
			} else if amount > 0 {
				l.Infof("Installed %d packages", len(conf.Packages))
			} else {
				l.Infof("Did nothing, exiting")
			}
		} else {
			l.Infof("No temporary packages found, skipped removing them")
		}

		l.Info("Installing permanent packages")
		err, amount := manager.Install(conf.Packages)
		if err != nil {
			l.Errorf("Failed to install packages: %s", err)
		} else if amount > 0 {
			l.Infof("Installed %d packages", len(conf.Packages))
		} else {
			l.Infof("Did nothing, exiting")
		}

		// TODO: execute configuration management here
	},
}
