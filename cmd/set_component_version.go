package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"k8s-cluster-upgrade-tool/config"
	"k8s-cluster-upgrade-tool/internal/api/k8s"
	"os/exec"
	"strings"
)

var setComponentVersionCmd = &cobra.Command{
	Use:   "setComponentVersion",
	Short: "Sets the value of a component running in the cluster to the passed value",
	Long: `Sets the value of a component running in the cluster to the passed value,
as of now will support setting the value for aws-node, cluster-autoscaler, kube-proxy, coredns
Usage:
$ k8s-cluster-upgrade-tool setComponentVersion valid-cluster-name aws-node my-version`,
	Args: cobra.ExactArgs(3),
	PreRun: func(cmd *cobra.Command, args []string) {
		err := config.Configuration.ValidatePassedComponentVersions(args[1], args[2])
		if err != nil {
			log.Fatalf("%s", err)
		}
		if config.Configuration.IsClusterNameValid(args[0]) {
			fmt.Println("Setting kubernetes context to", args[0])
			setK8sContext(args[0])
		} else {
			log.Fatal("Please pass a valid clusterName")
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		componentName, imageTag := args[1], args[2]
		switch componentName {
		case "coredns":
			setComponentVersion(imageTag, componentName, "deployment.apps/coredns", "deployment")
		case "kube-proxy":
			setComponentVersion(imageTag, componentName, "daemonset.apps/kube-proxy", "daemonset")
		case "aws-node":
			setComponentVersion(imageTag, componentName, "daemonset.apps/aws-node", "daemonset")
		case "cluster-autoscaler":
			log.Printf("for %s, please update the component via helm as we maintain the charts for the same.", componentName)
		}
	},
}

func init() {
	RootCmd.AddCommand(setComponentVersionCmd)

	// TODO Move the flags to required ones similar to taint-and-drain-asg command
}

func setComponentVersion(imageTag string, componentName string, k8sSetQueryCmdObject string, componentK8sObject string) {
	// get current imagePrefix
	args := strings.Fields(k8s.KubectlGetImageCommand(componentK8sObject, componentName))
	output, err := exec.Command(args[0], args[1:]...).Output()
	if err != nil {
		log.Fatal("There was an error while fetching the image of the component from the cluster: ", err)
	}

	imagePrefix, err := k8s.ParseComponentImage(string(output), "imagePrefix")
	if err != nil {
		log.Fatal("There was an error while parsing the image prefix step: ", err)
	}
	containerImage := imagePrefix + ":" + imageTag

	args = strings.Fields(k8s.KubectlSetImageCommand(k8sSetQueryCmdObject, componentName, containerImage))
	cmd := exec.Command(args[0], args[1:]...)
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%s has been set to %s in cluster \n", componentName, imageTag)
}
