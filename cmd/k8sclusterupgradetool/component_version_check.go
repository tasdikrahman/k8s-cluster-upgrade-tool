package k8sclusterupgradetool

import (
	"github.com/deliveryhero/k8s-cluster-upgrade-tool/config"
	"github.com/deliveryhero/k8s-cluster-upgrade-tool/internal/api/k8s"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os/exec"
	"strings"
)

func init() {
	componentVersionCmd.AddCommand(postUpgradeCheckCmd)

	postUpgradeCheckCmd.Flags().StringP("cluster", "c", "",
		"Example cluster name input valid-cluster-name, check with team for a full list of valid clusters")
	//nolint
	postUpgradeCheckCmd.MarkFlagRequired("cluster")
}

var postUpgradeCheckCmd = &cobra.Command{
	Use:   "check",
	Short: "Runs post upgrade checks on a cluster",
	Long: `Just checks for a cluster to see whether all the components have been upgraded or not
Usage:
$ k8sclusterupgradetool component version check -c=valid-cluster-name`,
	Run: func(cmd *cobra.Command, args []string) {
		cluster, _ := cmd.Flags().GetString("cluster")
		// Read config from file
		configFileName, configFileType, configFilePath := config.FileMetadata()
		configuration, err := config.Read(configFileName, configFileType, configFilePath)
		if err != nil {
			log.Fatalln(err)
		}

		log.Println("Config file used:", viper.ConfigFileUsed())
		log.Printf("aws-node version read from config: %s\n", viper.Get("components.aws-node"))
		log.Printf("coredns version read from config: %s", viper.Get("components.coredns"))
		log.Printf("kube-proxy version read from config: %s", viper.Get("components.kube-proxy"))
		log.Printf("cluster-autoscaler version read from config: %s", viper.Get("components.cluster-autoscaler"))

		if configuration.IsClusterNameValid(cluster) {
			log.Println("Setting kubernetes context to", cluster)
			k8s.SetK8sContext(cluster)
		} else {
			log.Fatal("Please pass a valid clusterName")
		}

		log.Println("running post upgrade checks")
		checkAwsNodeComponentVersion(cluster, configuration)
		checkKubeProxyComponentVersion(cluster, configuration)
		checkCoreDnsComponentVersion(cluster, configuration)
		checkClusterAutoscalerVersion(cluster, configuration)
	},
}

func checkAwsNodeComponentVersion(clusterName string, configuration config.Configurations) {
	log.Println("Checking aws-node version")
	// TODO: Change this to use to k8s client-go
	k8sObject, err := configuration.GetK8sObjectForCluster(clusterName, "aws-node")
	if err != nil {
		log.Fatalln("Error: there was an error while retrieving the k8sobject name and object type from the config")
	}
	args := strings.Fields(k8s.KubectlGetImageCommand(k8sObject.ObjectType, k8sObject.DeploymentName, k8sObject.Namespace))

	output, err := exec.Command(args[0], args[1:]...).Output()
	if err != nil {
		log.Fatalln("Error: there was an issue while retrieving the information from the cluster for the aws-node component")
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
	k8sObject, err := configuration.GetK8sObjectForCluster(clusterName, "kube-proxy")
	if err != nil {
		log.Fatalln("Error: there was an error while retrieving the k8sobject name and object type from the config")
	}
	args := strings.Fields(k8s.KubectlGetImageCommand(k8sObject.ObjectType, k8sObject.DeploymentName, k8sObject.Namespace))

	output, err := exec.Command(args[0], args[1:]...).Output()
	if err != nil {
		log.Fatalln("Error: there was an issue while retrieving the information from the cluster for the kube-proxy component")
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
	k8sObject, err := configuration.GetK8sObjectForCluster(clusterName, "coredns")
	if err != nil {
		log.Fatalln("Error: there was an error while retrieving the k8sobject name and object type from the config")
	}
	args := strings.Fields(k8s.KubectlGetImageCommand(k8sObject.ObjectType, k8sObject.DeploymentName, k8sObject.Namespace))
	output, err := exec.Command(args[0], args[1:]...).Output()
	if err != nil {
		log.Fatalln("Error: there was an issue while retrieving the information from the cluster for the coredns component")
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
	k8sObject, err := configuration.GetK8sObjectForCluster(clusterName, "cluster-autoscaler")
	if err != nil {
		log.Fatalln("Error: there was an error while retrieving the k8sobject name and object type from the config")
	}
	args := strings.Fields(k8s.KubectlGetImageCommand(k8sObject.ObjectType, k8sObject.DeploymentName, k8sObject.Namespace))
	output, err := exec.Command(args[0], args[1:]...).Output()
	if err != nil {
		log.Fatalln("Error: there was an issue while retrieving the information from the cluster for the cluster-autoscaler component")
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
