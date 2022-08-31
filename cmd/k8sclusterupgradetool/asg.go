package k8sclusterupgradetool

import (
	"fmt"
	"github.com/spf13/cobra"
)

var asgCmd = &cobra.Command{
	Use: "asg",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Autoscaling Group Operations")
		fmt.Println("Run 'k8sclusterupgradetool asg --help' to see the available commands")
	},
}

func init() {
	RootCmd.AddCommand(asgCmd)
}
