package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
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
		if config.Configuration.IsClusterNameValid(args[0]) {
			fmt.Println("Setting kubernetes context to", args[0])
			setK8sContext(args[0])
		} else {
			log.Fatal("Please pass a valid clusterName")
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("running post upgrade checks")
		checkAwsNodeComponentVersion()
		checkKubeProxyComponentVersion()
		checkCoreDnsComponentVersion()
		checkClusterAutoscalerVersion(args[0])
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

func checkAwsNodeComponentVersion() {
	fmt.Println("Checking aws-node version")
	// TODO: Change this to use to k8s client-go
	args := strings.Fields(k8s.KubectlGetImageCommand("daemonset", "aws-node"))
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

func checkKubeProxyComponentVersion() {
	fmt.Println("Checking kube-proxy version")
	// TODO: Change this to use to k8s client-go
	args := strings.Fields(k8s.KubectlGetImageCommand("daemonset", "kube-proxy"))
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

func checkCoreDnsComponentVersion() {
	fmt.Println("Checking core-dns version")
	// TODO: Change this to use to k8s client-go
	args := strings.Fields(k8s.KubectlGetImageCommand("deployment", "coredns"))
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

func checkClusterAutoscalerVersion(clusterName string) {
	fmt.Println("Checking cluster-autoscaler version")
	// TODO: Change this to use to k8s client-go
	deploymentName, err := k8s.GetClusterAutoscalerDeploymentNameForCluster(clusterName)
	if err != nil {
		log.Fatal(err)
	}
	args := strings.Fields(k8s.KubectlGetImageCommand("deployment", deploymentName))
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
