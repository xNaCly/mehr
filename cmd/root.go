package cmd

import (
	"fmt"
	"os"

	"github.com/xnacly/mehr/config"

	"github.com/spf13/cobra"
	l "github.com/xnacly/mehr/log"
)

var rootCmd = &cobra.Command{
	Use:     "mehr",
	Version: "0.0.1-dev",
	Short:   "Mehr is a operating system-independent package and system configuration manager",
	Long: `A Fast and Flexible package and configuration manager build in Go. 
Full configuration is available at https://github.com/xnacly/mehr`,
	Run: func(cmd *cobra.Command, args []string) {
		configPath := config.LookUp()
		conf, err := config.Get(configPath)

		if err != nil {
			l.Warn(err)
		}

		l.Infof("Got %q for Configuration.PackageManager", conf.PackageManager)
	},
}

func init() {
	rootCmd.PersistentFlags().BoolP("force", "f", false, "skip checks, may break system configuration")
}

func Root() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
