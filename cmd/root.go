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
		force, _ := cmd.Flags().GetBool("force")
		if config.IsRoot() {
			msg := "mehr started with elevated privileges, this can damage your system and requires a separate configuration,"
			if !force {
				l.Error(msg, "use --force to omit this check, exiting")
				return
			}
			l.Warn(msg, "got --force, not exiting")
		}
		configPath := config.LookUp()
		_, err := config.Get(configPath)

		if err != nil {
			l.Warn(err)
		}

		cmd.Help()
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
