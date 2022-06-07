package config

import (
	"errors"
	"github.com/spf13/viper"
	"log"
)

const (
	FileName = "config"
	FileType = "yaml"
	FilePath = "$HOME/.k8s-cluster-upgrade-tool"
)

type Configurations struct {
	Components  ComponentVersionConfigurations `mapstructure:"components"`
	ClusterList []ClusterListConfiguration     `mapstructure:"clusterlist"`
}

// reference: https://stackoverflow.com/questions/63889004/how-to-access-specific-items-in-an-array-from-viper
type ClusterListConfiguration struct {
	Name                    string    `mapstructure:"Name"`
	AwsRegion               string    `mapstructure:"AwsRegion"`
	AwsAccount              string    `mapstructure:"AwsAccount"`
	AwsNodeObject           K8sObject `mapstructure:"AwsNodeObject"`
	ClusterAutoscalerObject K8sObject `mapstructure:"ClusterAutoscalerObject"`
	CoreDnsObject           K8sObject `mapstructure:"CoreDnsObject"`
	KubeProxyObject         K8sObject `mapstructure:"KubeProxyObject"`
}

type K8sObject struct {
	Name string `mapstructure:"name"`
	Type string `mapstructure:"type"`
}

type ComponentVersionConfigurations struct {
	AwsNode           string `mapstructure:"aws-node"`
	ClusterAutoscaler string `mapstructure:"cluster-autoscaler"`
	CoreDns           string `mapstructure:"coredns"`
	KubeProxy         string `mapstructure:"kube-proxy"`
}

// TODO Add spec to check for unique cluster name list in config
func (c Configurations) IsClusterListConfigurationValid() bool {
	valid := true
	for _, cluster := range c.ClusterList {
		if cluster.Name == "" || cluster.AwsRegion == "" || cluster.AwsAccount == "" || cluster.AwsNodeObject.Name == "" || cluster.AwsNodeObject.Type == "" || cluster.ClusterAutoscalerObject.Name == "" || cluster.ClusterAutoscalerObject.Type == "" || cluster.CoreDnsObject.Name == "" || cluster.CoreDnsObject.Type == "" || cluster.KubeProxyObject.Name == "" || cluster.KubeProxyObject.Type == "" {
			valid = false
		}
	}
	return valid
}

func (c Configurations) IsComponentVersionConfigurationsValid() bool {
	valid := true
	if c.Components.CoreDns == "" || c.Components.AwsNode == "" || c.Components.ClusterAutoscaler == "" || c.Components.KubeProxy == "" {
		valid = false
	}
	return valid
}

func (c Configurations) IsClusterNameValid(clusterName string) bool {
	contains := false
	for _, cluster := range c.ClusterList {
		if cluster.Name == clusterName {
			contains = true
		}
	}
	return contains
}

func (c Configurations) GetK8sObjectNameAndObjectTypeForCluster(clusterName, k8sObject string) (objectName, objectType string, err error) {
	for _, cluster := range c.ClusterList {
		if cluster.Name == clusterName {
			switch k8sObject {
			case "aws-node":
				return cluster.AwsNodeObject.Name, cluster.AwsNodeObject.Type, nil
			case "cluster-autoscaler":
				return cluster.ClusterAutoscalerObject.Name, cluster.ClusterAutoscalerObject.Type, nil
			case "kube-proxy":
				return cluster.KubeProxyObject.Name, cluster.KubeProxyObject.Type, nil
			case "coredns":
				return cluster.CoreDnsObject.Name, cluster.CoreDnsObject.Type, nil
			default:
				return "", "", errors.New("please pass any of the components between aws-node, coredns, cluster-autoscaler, kube-proxy")
			}
		}
	}
	return "", "", errors.New("please check if you passed a valid cluster name")
}

func (c Configurations) GetAwsAccountAndRegionForCluster(clusterName string) (awsAccount, awsRegion string, err error) {
	for _, cluster := range c.ClusterList {
		if cluster.Name == clusterName {
			return cluster.AwsAccount, cluster.AwsRegion, nil
		}
	}
	return "", "", errors.New("no awsAccount and awsRegion was found for the passed clusterName")
}

func (c Configurations) ValidatePassedComponentVersions(componentName, componentVersion string) error {
	switch componentName {
	case "aws-node":
		if componentVersion != c.Components.AwsNode {
			return errors.New("aws-node component version passed doesn't match the version in config, please check the value in config file")
		}
	case "cluster-autoscaler":
		if componentVersion != c.Components.ClusterAutoscaler {
			return errors.New("cluster-autoscaler component version passed doesn't match the version in config, please check the value in config file")
		}
	case "kube-proxy":
		if componentVersion != c.Components.KubeProxy {
			return errors.New("kube-proxy component version passed doesn't match the version in config, please check the value in config file")
		}
	case "coredns":
		if componentVersion != c.Components.CoreDns {
			return errors.New("coredns component version passed doesn't match the version in config, please check the value in config file")
		}
	default:
		return errors.New("please pass a valid component name from this list [coredns, cluster-autoscaler, kube-proxy, aws-node]")
	}

	return nil
}

func Read(fileName, fileType, filePath string) (config Configurations, err error) {
	viper.SetConfigName(fileName)
	viper.SetConfigType(fileType)
	viper.AddConfigPath(filePath)
	err = viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return Configurations{}, errors.New("error finding config file. Does it exist? Please create it in $HOME/.k8s-cluster-upgrade-tool/config.yaml if not")
		} else {
			return Configurations{}, errors.New("error reading from config file")
		}
	}
	log.Println("Config file used:", viper.ConfigFileUsed())

	err = viper.Unmarshal(&config)
	if err != nil {
		return Configurations{}, errors.New("error un marshaling config file")
	}

	// check for the mandatory config file variables being read
	if !config.IsComponentVersionConfigurationsValid() {
		return Configurations{}, errors.New("mandatory component version of either aws-node, coredns, kube-proxy or cluster-autoscaler not set in config file")
	}

	if !config.IsClusterListConfigurationValid() {
		return Configurations{}, errors.New("one of the clusterlist elements has either Name, AwsRegion, AwsAccount, AwsNodeObject, ClusterAutoscalerObject, KubeProxyObject, CoreDnsObject is missing")
	}

	log.Printf("aws-node version read from config: %s\n", viper.Get("components.aws-node"))
	log.Printf("coredns version read from config: %s", viper.Get("components.coredns"))
	log.Printf("kube-proxy version read from config: %s", viper.Get("components.kube-proxy"))
	log.Printf("cluster-autoscaler version read from config: %s", viper.Get("components.cluster-autoscaler"))
	return config, nil
}

func FileMetadata() (fileName, filePath, fileType string) {
	return FileName, FileType, FilePath
}
