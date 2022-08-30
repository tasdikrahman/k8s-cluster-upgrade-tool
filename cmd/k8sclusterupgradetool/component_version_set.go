package k8sclusterupgradetool

import (
	"fmt"
	"github.com/deliveryhero/k8s-cluster-upgrade-tool/config"
	"github.com/deliveryhero/k8s-cluster-upgrade-tool/internal/api/k8s"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os/exec"
	"strings"
)

var setComponentVersionCmd = &cobra.Command{
	Use:   "set",
	Short: "Sets the value of a component running in the cluster to the passed value",
	Long: `Sets the value of a component running in the cluster to the passed value,
as of now will support setting the value for aws-node, cluster-autoscaler, kube-proxy, coredns
Usage:
$ k8sclusterupgradetool component version set -c=valid-cluster-name -o=aws-node -v=my-version`,
	Run: func(cmd *cobra.Command, args []string) {
		// Parse flag values
		cluster, _ := cmd.Flags().GetString("cluster")
		k8sComponent, _ := cmd.Flags().GetString("component-object")
		k8sComponentVersion, _ := cmd.Flags().GetString("component-object-version")

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

		err = configuration.ValidatePassedComponentVersions(k8sComponent, k8sComponentVersion)
		if err != nil {
			log.Fatalf("%s", err)
		}

		if configuration.IsClusterNameValid(cluster) {
			log.Println("Setting kubernetes context to", cluster)
			k8s.SetK8sContext(cluster)
		} else {
			log.Fatal("Please pass a valid clusterName")
		}

		componentName, imageTag := k8sComponent, k8sComponentVersion
		switch componentName {
		case "coredns":
			k8sObject, err := configuration.GetK8sObjectForCluster(cluster, "coredns")
			if err != nil {
				log.Fatalln("There was an error reading config from the config file")
			}
			setComponentVersion(imageTag, componentName, fmt.Sprintf("%s.apps/%s", k8sObject.ObjectType, k8sObject.DeploymentName), k8sObject.ObjectType, k8sObject.ContainerName, k8sObject.Namespace)
		case "kube-proxy":
			k8sObject, err := configuration.GetK8sObjectForCluster(cluster, "kube-proxy")
			if err != nil {
				log.Println(err)
			}
			setComponentVersion(imageTag, componentName, fmt.Sprintf("%s.apps/%s", k8sObject.ObjectType, k8sObject.DeploymentName), k8sObject.ObjectType, k8sObject.ContainerName, k8sObject.Namespace)
		case "aws-node":
			k8sObject, err := configuration.GetK8sObjectForCluster(cluster, "aws-node")
			if err != nil {
				log.Println(err)
			}
			setComponentVersion(imageTag, componentName, fmt.Sprintf("%s.apps/%s", k8sObject.ObjectType, k8sObject.DeploymentName), k8sObject.ObjectType, k8sObject.ContainerName, k8sObject.Namespace)
		case "cluster-autoscaler":
			k8sObject, err := configuration.GetK8sObjectForCluster(cluster, "cluster-autoscaler")
			if err != nil {
				log.Println(err)
			}
			setComponentVersion(imageTag, componentName, fmt.Sprintf("%s.apps/%s", k8sObject.ObjectType, k8sObject.DeploymentName), k8sObject.ObjectType, k8sObject.ContainerName, k8sObject.Namespace)
		default:
			log.Println("please check the passed components, the supported components are cluster-autoscaler, kube-proxy, coredns, aws-node")
		}
	},
}

func init() {
	componentVersionCmd.AddCommand(setComponentVersionCmd)

	setComponentVersionCmd.Flags().StringP("cluster", "c", "",
		"Example cluster name input valid-cluster-name, check with team for a full list of valid clusters")
	setComponentVersionCmd.Flags().StringP("component-object", "o", "",
		"K8s cluster component being set, currently supported ones: eg: aws-node, cluster-autoscaler, kube-proxy, coredns")
	setComponentVersionCmd.Flags().StringP("component-object-version", "v", "",
		"k8s component version to be set for the k8s component, currently supported ones: eg: aws-node, cluster-autoscaler, kube-proxy, coredns")
	//nolint
	nodeTaintAndDrainCmd.MarkFlagRequired("cluster")
	//nolint
	nodeTaintAndDrainCmd.MarkFlagRequired("component-object")
	//nolint
	nodeTaintAndDrainCmd.MarkFlagRequired("component-object-version")
}

func setComponentVersion(imageTag, componentName, k8sSetQueryCmdObject, componentK8sObject, containerName, namespace string) {
	// get current imagePrefix
	args := strings.Fields(k8s.KubectlGetImageCommand(componentK8sObject, componentName, namespace))
	output, err := exec.Command(args[0], args[1:]...).Output()
	if err != nil {
		log.Fatalln("There was an error while fetching the image of the component from the cluster: ", err)
	}

	imagePrefix, err := k8s.ParseComponentImage(string(output), "imagePrefix")
	if err != nil {
		log.Fatalln("There was an error while parsing the image prefix step: ", err)
	}
	containerImage := imagePrefix + ":" + imageTag

	args = strings.Fields(k8s.KubectlSetImageCommand(k8sSetQueryCmdObject, containerName, containerImage, namespace))
	cmd := exec.Command(args[0], args[1:]...)
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%s has been set to %s in cluster \n", componentName, imageTag)
}
