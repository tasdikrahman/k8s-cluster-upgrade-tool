package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"os/exec"
	"strings"

	"k8s-cluster-upgrade-tool/config"
	"k8s-cluster-upgrade-tool/internal/api/k8s"
)

var postUpgradeCheckCmd = &cobra.Command{
	Use:   "postUpgradeCheck",
	Short: "Runs post upgrade checks on a cluster",
	Long: `Just checks for a cluster to see whether all the components have been upgraded or not
Usage:
$ k8s-cluster-upgrade-tool postUpgradeCheck valid-cluster-name`,
	Args: cobra.ExactArgs(1),
	PreRun: func(cmd *cobra.Command, args []string) {
		// Read config from file
		configFileName, configFileType, configFilePath := config.FileMetadata()
		configuration, err := config.Read(configFileName, configFileType, configFilePath)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if configuration.IsClusterNameValid(args[0]) {
			fmt.Println("Setting kubernetes context to", args[0])
			setK8sContext(args[0])
		} else {
			log.Fatal("Please pass a valid clusterName")
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		// Read config from file
		configFileName, configFileType, configFilePath := config.FileMetadata()
		configuration, err := config.Read(configFileName, configFileType, configFilePath)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println("running post upgrade checks")
		checkAwsNodeComponentVersion(args[0], configuration)
		checkKubeProxyComponentVersion(args[0], configuration)
		checkCoreDnsComponentVersion(args[0], configuration)
		checkClusterAutoscalerVersion(args[0], configuration)
	},
}

func init() {
	RootCmd.AddCommand(postUpgradeCheckCmd)

	// TODO Move the flags to required ones similar to taint-and-drain-asg command
}

func setK8sContext(clusterName string) {
	command := "kubectl"
	arg01 := "config"
	arg02 := "use-context"

	// TODO: change this to use client-go
	cmd := exec.Command(command, arg01, arg02, clusterName)
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func checkAwsNodeComponentVersion(clusterName string, configuration config.Configurations) {
	fmt.Println("Checking aws-node version")
	// TODO: Change this to use to k8s client-go
	k8sObjectName, k8sObjectType, err := configuration.GetK8sObjectNameAndObjectTypeForCluster(clusterName, "aws-node")
	if err != nil {
		log.Fatal(err)
	}
	args := strings.Fields(k8s.KubectlGetImageCommand(k8sObjectType, k8sObjectName))
	output, err := exec.Command(args[0], args[1:]...).Output()
	if err != nil {
		log.Fatal(err)
	}

	imageTag, err := k8s.ParseComponentImage(string(output), "imageTag")
	if err != nil {
		log.Fatal(err)
	}
	if imageTag == viper.Get("components.aws-node") {
		fmt.Printf("AWS Node Version on %s ✓ \n", viper.Get("components.aws-node"))
	} else {
		fmt.Printf("aws-node needs to be updated, is currently on %s, desired version: %s\n", imageTag,
			viper.Get("components.aws-node"))
	}
}

func checkKubeProxyComponentVersion(clusterName string, configuration config.Configurations) {
	fmt.Println("Checking kube-proxy version")
	// TODO: Change this to use to k8s client-go
	k8sObjectName, k8sObjectType, err := configuration.GetK8sObjectNameAndObjectTypeForCluster(clusterName, "kube-proxy")
	if err != nil {
		log.Fatal(err)
	}
	args := strings.Fields(k8s.KubectlGetImageCommand(k8sObjectType, k8sObjectName))
	output, err := exec.Command(args[0], args[1:]...).Output()
	if err != nil {
		log.Fatal(err)
	}

	imageTag, err := k8s.ParseComponentImage(string(output), "imageTag")
	if err != nil {
		log.Fatal(err)
	}

	if imageTag == viper.Get("components.kube-proxy") {
		fmt.Printf("kube-proxy on %s ✓ \n", viper.Get("components.kube-proxy"))
	} else {
		fmt.Printf("kube-proxy needs to be updated, is currently on %s, desired version: %s\n", imageTag,
			viper.Get("components.kube-proxy"))
	}
}

func checkCoreDnsComponentVersion(clusterName string, configuration config.Configurations) {
	fmt.Println("Checking coredns version")
	// TODO: Change this to use to k8s client-go
	k8sObjectName, k8sObjectType, err := configuration.GetK8sObjectNameAndObjectTypeForCluster(clusterName, "coredns")
	if err != nil {
		log.Fatal(err)
	}
	args := strings.Fields(k8s.KubectlGetImageCommand(k8sObjectType, k8sObjectName))
	output, err := exec.Command(args[0], args[1:]...).Output()
	if err != nil {
		log.Fatal(err)
	}

	imageTag, err := k8s.ParseComponentImage(string(output), "imageTag")
	if err != nil {
		log.Fatal(err)
	}

	if imageTag == viper.Get("components.coredns") {
		fmt.Printf("core-dns on %s ✓ \n", viper.Get("components.coredns"))
	} else {
		fmt.Printf("core-dns needs to be updated, is currently on %s, desired version: %s\n", imageTag,
			viper.Get("components.coredns"))
	}
}

func checkClusterAutoscalerVersion(clusterName string, configuration config.Configurations) {
	fmt.Println("Checking cluster-autoscaler version")
	// TODO: Change this to use to k8s client-go
	k8sObjectName, k8sObjectType, err := configuration.GetK8sObjectNameAndObjectTypeForCluster(clusterName, "cluster-autoscaler")
	if err != nil {
		log.Fatal(err)
	}
	args := strings.Fields(k8s.KubectlGetImageCommand(k8sObjectType, k8sObjectName))
	output, err := exec.Command(args[0], args[1:]...).Output()
	if err != nil {
		log.Fatal(err)
	}

	imageTag, err := k8s.ParseComponentImage(string(output), "imageTag")
	if err != nil {
		log.Fatal(err)
	}

	if imageTag == viper.Get("components.cluster-autoscaler") {
		fmt.Printf("cluster-autoscaler on %s ✓ \n", viper.Get("components.cluster-autoscaler"))
	} else {
		fmt.Printf("cluster-autoscaler needs to be updated, is currently on %s, desired version: %s\n", imageTag,
			viper.Get("components.cluster-autoscaler"))
	}
}
