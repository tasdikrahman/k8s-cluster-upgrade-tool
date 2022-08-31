package k8sclusterupgradetool

import (
	"fmt"
	"github.com/spf13/cobra"
)

var componentVersionCmd = &cobra.Command{
	Use: "version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("K8s cluster component version operations")
		fmt.Println("Run 'k8sclusterupgradetool component version --help' to see the available commands")
	},
}

func init() {
	componentCmd.AddCommand(componentVersionCmd)
}
