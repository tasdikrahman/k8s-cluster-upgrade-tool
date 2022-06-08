package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s-cluster-upgrade-tool/config"
	"k8s-cluster-upgrade-tool/internal/api/k8s"
	"log"
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
		// Read config from file
		configFileName, configFileType, configFilePath := config.FileMetadata()
		configuration, err := config.Read(configFileName, configFileType, configFilePath)
		if err != nil {
			log.Fatalln("There was an error reading config from the config file")
		}

		log.Println("Config file used:", viper.ConfigFileUsed())
		log.Printf("aws-node version read from config: %s\n", viper.Get("components.aws-node"))
		log.Printf("coredns version read from config: %s", viper.Get("components.coredns"))
		log.Printf("kube-proxy version read from config: %s", viper.Get("components.kube-proxy"))
		log.Printf("cluster-autoscaler version read from config: %s", viper.Get("components.cluster-autoscaler"))

		err = configuration.ValidatePassedComponentVersions(args[1], args[2])
		if err != nil {
			log.Fatalf("%s", err)
		}

		if configuration.IsClusterNameValid(args[0]) {
			log.Println("Setting kubernetes context to", args[0])
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
			log.Fatal(err)
		}

		componentName, imageTag := args[1], args[2]
		switch componentName {
		case "coredns":
			k8sObjectName, k8sObjectType, err := configuration.GetK8sObjectNameAndObjectTypeForCluster(args[0], "coredns")
			if err != nil {
				log.Fatalln("There was an error reading config from the config file")
			}
			setComponentVersion(imageTag, componentName, fmt.Sprintf("%s.apps/%s", k8sObjectType, k8sObjectName), k8sObjectType)
		case "kube-proxy":
			k8sObjectName, k8sObjectType, err := configuration.GetK8sObjectNameAndObjectTypeForCluster(args[0], "kube-proxy")
			if err != nil {
				log.Println(err)
			}
			setComponentVersion(imageTag, componentName, fmt.Sprintf("%s.apps/%s", k8sObjectType, k8sObjectName), k8sObjectType)
		case "aws-node":
			k8sObjectName, k8sObjectType, err := configuration.GetK8sObjectNameAndObjectTypeForCluster(args[0], "aws-node")
			if err != nil {
				log.Println(err)
			}
			setComponentVersion(imageTag, componentName, fmt.Sprintf("%s.apps/%s", k8sObjectType, k8sObjectName), k8sObjectType)
		case "cluster-autoscaler":
			k8sObjectName, k8sObjectType, err := configuration.GetK8sObjectNameAndObjectTypeForCluster(args[0], "cluster-autoscaler")
			if err != nil {
				log.Println(err)
			}
			setComponentVersion(imageTag, componentName, fmt.Sprintf("%s.apps/%s", k8sObjectType, k8sObjectName), k8sObjectType)
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
		log.Fatalln("There was an error while fetching the image of the component from the cluster: ", err)
	}

	imagePrefix, err := k8s.ParseComponentImage(string(output), "imagePrefix")
	if err != nil {
		log.Fatalln("There was an error while parsing the image prefix step: ", err)
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
