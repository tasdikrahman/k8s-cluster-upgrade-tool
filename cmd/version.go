package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

const (
	K8sClusterUpgradeToolVersion = "v0.1.0"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Shows the version of the release of k8s-cluster-upgrade-tool binary",
	Args:  cobra.MaximumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Version: ", K8sClusterUpgradeToolVersion)
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
