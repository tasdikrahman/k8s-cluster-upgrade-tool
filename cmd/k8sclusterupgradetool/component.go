package k8sclusterupgradetool

import (
	"fmt"
	"github.com/spf13/cobra"
)

var componentCmd = &cobra.Command{
	Use: "component",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("K8s cluster component operations")
		fmt.Println("Run 'k8sclusterupgradetool component --help' to see the available commands")
	},
}

func init() {
	RootCmd.AddCommand(componentCmd)
}
