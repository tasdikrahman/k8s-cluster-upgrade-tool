package k8sclusterupgradetool

import (
	"fmt"
	"github.com/deliveryhero/k8s-cluster-upgrade-tool/config"
	"github.com/deliveryhero/k8s-cluster-upgrade-tool/internal/api/k8s"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/client-go/kubernetes"
	"log"
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
			log.Println("running post upgrade checks")
			k8sClient, err := k8s.KubeClientInit(cluster)
			if err != nil {
				log.Fatal("There was an error initializing the k8sclient with the passed cluster context")
			}

			err = checkComponentVersion("aws-node", cluster, configuration, k8sClient)
			if err != nil {
				log.Fatalf("error while checking for aws-node component version: %v", err)
			}
			err = checkComponentVersion("kube-proxy", cluster, configuration, k8sClient)
			if err != nil {
				log.Fatalf("error while checking for kube-proxy component version: %v", err)
			}
			err = checkComponentVersion("coredns", cluster, configuration, k8sClient)
			if err != nil {
				log.Fatalf("error while checking for coredns component version: %v", err)
			}
			err = checkComponentVersion("cluster-autoscaler", cluster, configuration, k8sClient)
			if err != nil {
				log.Fatalf("error while checking for cluster-autoscaler component version: %v", err)
			}
		} else {
			log.Fatal("Please pass a valid clusterName")
		}
	},
}

func checkComponentVersion(componentName, clusterName string, configuration config.Configurations, k8sClient kubernetes.Interface) error {
	log.Printf("Checking %s version\n", componentName)
	k8sObject, err := configuration.GetK8sObjectForCluster(clusterName, componentName)
	if err != nil {
		return err
	}

	containerImage, err := k8s.GetContainerImageForK8sObject(k8sClient, k8sObject.DeploymentName, k8sObject.ObjectType, k8sObject.Namespace)
	if err != nil {
		return err
	}

	imageTag, err := k8s.ParseComponentImage(containerImage, "imageTag")
	if err != nil {
		return err
	}

	viperQuery := fmt.Sprintf("components.%s", componentName)
	if imageTag == viper.Get(viperQuery) {
		log.Printf("%s Version on %s ✓ \n", componentName, viper.Get(viperQuery))
	} else {
		log.Printf("%s needs to be updated, is currently on %s, desired version: %s\n", componentName, imageTag, viper.Get(viperQuery))
	}
	return nil
}
