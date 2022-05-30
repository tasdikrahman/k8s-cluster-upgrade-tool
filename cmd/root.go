package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"k8s-cluster-upgrade-tool/config"
	"os"
)

var RootCmd = &cobra.Command{
	Use:   "k8s-cluster-upgrade-tool",
	Short: "k8s-cluster-upgrade-tool",
}

func Execute() {
	// Read config from file
	configFileName, configFileType, configFilePath := config.FileMetadata()
	err := config.Read(configFileName, configFileType, configFilePath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
