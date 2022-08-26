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
$ k8sclusterupgradetool component version set valid-cluster-name aws-node my-version`,
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

		componentName, imageTag := args[1], args[2]
		switch componentName {
		case "coredns":
			k8sObject, err := configuration.GetK8sObjectForCluster(args[0], "coredns")
			if err != nil {
				log.Fatalln("There was an error reading config from the config file")
			}
			setComponentVersion(imageTag, componentName, fmt.Sprintf("%s.apps/%s", k8sObject.ObjectType, k8sObject.DeploymentName), k8sObject.ObjectType, k8sObject.ContainerName, k8sObject.Namespace)
		case "kube-proxy":
			k8sObject, err := configuration.GetK8sObjectForCluster(args[0], "kube-proxy")
			if err != nil {
				log.Println(err)
			}
			setComponentVersion(imageTag, componentName, fmt.Sprintf("%s.apps/%s", k8sObject.ObjectType, k8sObject.DeploymentName), k8sObject.ObjectType, k8sObject.ContainerName, k8sObject.Namespace)
		case "aws-node":
			k8sObject, err := configuration.GetK8sObjectForCluster(args[0], "aws-node")
			if err != nil {
				log.Println(err)
			}
			setComponentVersion(imageTag, componentName, fmt.Sprintf("%s.apps/%s", k8sObject.ObjectType, k8sObject.DeploymentName), k8sObject.ObjectType, k8sObject.ContainerName, k8sObject.Namespace)
		case "cluster-autoscaler":
			k8sObject, err := configuration.GetK8sObjectForCluster(args[0], "cluster-autoscaler")
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

	// TODO Move the flags to required ones similar to taint-and-drain-asg command
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
