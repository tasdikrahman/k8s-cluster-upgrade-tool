package cmd

import (
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

		if configuration.IsClusterNameValid(args[0]) {
			log.Println("Setting kubernetes context to", args[0])
			k8s.SetK8sContext(args[0])
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

		log.Println("running post upgrade checks")
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

func checkAwsNodeComponentVersion(clusterName string, configuration config.Configurations) {
	log.Println("Checking aws-node version")
	// TODO: Change this to use to k8s client-go
	k8sObjectName, k8sObjectType, _, err := configuration.GetK8sObjectNameObjectTypeAndContainerNameForCluster(clusterName, "aws-node")
	if err != nil {
		log.Fatalln("Error: there was an error while retrieving the k8sobject name and object type from the config")
	}
	args := strings.Fields(k8s.KubectlGetImageCommand(k8sObjectType, k8sObjectName))
	output, err := exec.Command(args[0], args[1:]...).Output()
	if err != nil {
		log.Fatalln("Error: there was an issue while retrieving the information from the cluster for the component")
	}

	imageTag, err := k8s.ParseComponentImage(string(output), "imageTag")
	if err != nil {
		log.Fatalln("Error: there was an error parsing the image from the parsed command output")
	}
	if imageTag == viper.Get("components.aws-node") {
		log.Printf("AWS Node Version on %s ✓ \n", viper.Get("components.aws-node"))
	} else {
		log.Printf("aws-node needs to be updated, is currently on %s, desired version: %s\n", imageTag,
			viper.Get("components.aws-node"))
	}
}

func checkKubeProxyComponentVersion(clusterName string, configuration config.Configurations) {
	log.Println("Checking kube-proxy version")
	// TODO: Change this to use to k8s client-go
	k8sObjectName, k8sObjectType, _, err := configuration.GetK8sObjectNameObjectTypeAndContainerNameForCluster(clusterName, "kube-proxy")
	if err != nil {
		log.Fatalln("Error: there was an error while retrieving the k8sobject name and object type from the config")
	}
	args := strings.Fields(k8s.KubectlGetImageCommand(k8sObjectType, k8sObjectName))
	output, err := exec.Command(args[0], args[1:]...).Output()
	if err != nil {
		log.Fatalln("Error: there was an issue while retrieving the information from the cluster for the component")
	}

	imageTag, err := k8s.ParseComponentImage(string(output), "imageTag")
	if err != nil {
		log.Fatalln("Error: there was an error parsing the image from the parsed command output")
	}

	if imageTag == viper.Get("components.kube-proxy") {
		log.Printf("kube-proxy on %s ✓ \n", viper.Get("components.kube-proxy"))
	} else {
		log.Printf("kube-proxy needs to be updated, is currently on %s, desired version: %s\n", imageTag,
			viper.Get("components.kube-proxy"))
	}
}

func checkCoreDnsComponentVersion(clusterName string, configuration config.Configurations) {
	log.Println("Checking coredns version")
	// TODO: Change this to use to k8s client-go
	k8sObjectName, k8sObjectType, _, err := configuration.GetK8sObjectNameObjectTypeAndContainerNameForCluster(clusterName, "coredns")
	if err != nil {
		log.Fatalln("Error: there was an error while retrieving the k8sobject name and object type from the config")
	}
	args := strings.Fields(k8s.KubectlGetImageCommand(k8sObjectType, k8sObjectName))
	output, err := exec.Command(args[0], args[1:]...).Output()
	if err != nil {
		log.Fatalln("Error: there was an issue while retrieving the information from the cluster for the component")
	}

	imageTag, err := k8s.ParseComponentImage(string(output), "imageTag")
	if err != nil {
		log.Fatalln("Error: there was an error parsing the image from the parsed command output")
	}

	if imageTag == viper.Get("components.coredns") {
		log.Printf("core-dns on %s ✓ \n", viper.Get("components.coredns"))
	} else {
		log.Printf("core-dns needs to be updated, is currently on %s, desired version: %s\n", imageTag,
			viper.Get("components.coredns"))
	}
}

func checkClusterAutoscalerVersion(clusterName string, configuration config.Configurations) {
	log.Println("Checking cluster-autoscaler version")
	// TODO: Change this to use to k8s client-go
	k8sObjectName, k8sObjectType, _, err := configuration.GetK8sObjectNameObjectTypeAndContainerNameForCluster(clusterName, "cluster-autoscaler")
	if err != nil {
		log.Fatalln("Error: there was an error while retrieving the k8sobject name and object type from the config")
	}
	args := strings.Fields(k8s.KubectlGetImageCommand(k8sObjectType, k8sObjectName))
	output, err := exec.Command(args[0], args[1:]...).Output()
	if err != nil {
		log.Fatalln("Error: there was an issue while retrieving the information from the cluster for the component")
	}

	imageTag, err := k8s.ParseComponentImage(string(output), "imageTag")
	if err != nil {
		log.Fatalln("Error: there was an error parsing the image from the parsed command output")
	}

	if imageTag == viper.Get("components.cluster-autoscaler") {
		log.Printf("cluster-autoscaler on %s ✓ \n", viper.Get("components.cluster-autoscaler"))
	} else {
		log.Printf("cluster-autoscaler needs to be updated, is currently on %s, desired version: %s\n", imageTag,
			viper.Get("components.cluster-autoscaler"))
	}
}
