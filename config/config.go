package config

import (
	"errors"
	"github.com/spf13/viper"
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
	ClusterName             string    `mapstructure:"ClusterName"`
	AwsRegion               string    `mapstructure:"AwsRegion"`
	AwsAccount              string    `mapstructure:"AwsAccount"`
	AwsNodeObject           K8sObject `mapstructure:"AwsNodeObject"`
	ClusterAutoscalerObject K8sObject `mapstructure:"ClusterAutoscalerObject"`
	CoreDnsObject           K8sObject `mapstructure:"CoreDnsObject"`
	KubeProxyObject         K8sObject `mapstructure:"KubeProxyObject"`
}

type K8sObject struct {
	DeploymentName string `mapstructure:"DeploymentName"`
	ObjectType     string `mapstructure:"ObjectType"`
	ContainerName  string `mapstructure:"ContainerName"`
	Namespace      string `mapstructure:"Namespace"`
}

type ComponentVersionConfigurations struct {
	AwsNode           string `mapstructure:"aws-node"`
	ClusterAutoscaler string `mapstructure:"cluster-autoscaler"`
	CoreDns           string `mapstructure:"coredns"`
	KubeProxy         string `mapstructure:"kube-proxy"`
}

func (c Configurations) IsClusterListConfigurationValid() bool {
	valid := true
	clusterNameMap := map[string]string{}
	for _, cluster := range c.ClusterList {
		if _, present := clusterNameMap[cluster.ClusterName]; present {
			valid = false
		}
		clusterNameMap[cluster.ClusterName] = "present"

		clusterName := cluster.ClusterName == ""
		awsRegion := cluster.AwsRegion == ""
		awsAccount := cluster.AwsAccount == ""
		awsNodeObject := cluster.AwsNodeObject.DeploymentName == "" || cluster.AwsNodeObject.ObjectType == "" || cluster.AwsNodeObject.ContainerName == "" || cluster.AwsNodeObject.Namespace == ""
		clusterAutoscaler := cluster.ClusterAutoscalerObject.DeploymentName == "" || cluster.ClusterAutoscalerObject.ObjectType == "" || cluster.ClusterAutoscalerObject.ContainerName == "" || cluster.ClusterAutoscalerObject.Namespace == ""
		coreDns := cluster.CoreDnsObject.DeploymentName == "" || cluster.CoreDnsObject.ObjectType == "" || cluster.CoreDnsObject.ContainerName == "" || cluster.CoreDnsObject.Namespace == ""
		kubeProxy := cluster.KubeProxyObject.DeploymentName == "" || cluster.KubeProxyObject.ObjectType == "" || cluster.KubeProxyObject.ContainerName == "" || cluster.KubeProxyObject.Namespace == ""

		if clusterName || awsRegion || awsAccount || awsNodeObject || clusterAutoscaler || coreDns || kubeProxy {
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
		if cluster.ClusterName == clusterName {
			contains = true
		}
	}
	return contains
}

func (c Configurations) GetK8sObjectForCluster(clusterName, k8sObjectType string) (k8sObject K8sObject, err error) {
	for _, cluster := range c.ClusterList {
		if cluster.ClusterName == clusterName {
			switch k8sObjectType {
			case "aws-node":
				return K8sObject{
					DeploymentName: cluster.AwsNodeObject.DeploymentName,
					ObjectType:     cluster.AwsNodeObject.ObjectType,
					ContainerName:  cluster.AwsNodeObject.ContainerName,
					Namespace:      cluster.AwsNodeObject.Namespace,
				}, nil
			case "cluster-autoscaler":
				return K8sObject{
					DeploymentName: cluster.ClusterAutoscalerObject.DeploymentName,
					ObjectType:     cluster.ClusterAutoscalerObject.ObjectType,
					ContainerName:  cluster.ClusterAutoscalerObject.ContainerName,
					Namespace:      cluster.ClusterAutoscalerObject.Namespace,
				}, nil
			case "kube-proxy":
				return K8sObject{
					DeploymentName: cluster.KubeProxyObject.DeploymentName,
					ObjectType:     cluster.KubeProxyObject.ObjectType,
					ContainerName:  cluster.KubeProxyObject.ContainerName,
					Namespace:      cluster.KubeProxyObject.Namespace,
				}, nil
			case "coredns":
				return K8sObject{
					DeploymentName: cluster.CoreDnsObject.DeploymentName,
					ObjectType:     cluster.CoreDnsObject.ObjectType,
					ContainerName:  cluster.CoreDnsObject.ContainerName,
					Namespace:      cluster.CoreDnsObject.Namespace,
				}, nil
			default:
				return K8sObject{
					DeploymentName: "",
					ObjectType:     "",
					ContainerName:  "",
					Namespace:      "",
				}, errors.New("please pass any of the components between aws-node, coredns, cluster-autoscaler, kube-proxy")
			}
		}
	}
	return K8sObject{
		DeploymentName: "",
		ObjectType:     "",
		ContainerName:  "",
		Namespace:      "",
	}, errors.New("please check if you passed a valid cluster name")
}

func (c Configurations) GetAwsAccountAndRegionForCluster(clusterName string) (awsAccount, awsRegion string, err error) {
	for _, cluster := range c.ClusterList {
		if cluster.ClusterName == clusterName {
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

	err = viper.Unmarshal(&config)
	if err != nil {
		return Configurations{}, errors.New("error un marshaling config file")
	}

	// check for the mandatory config file variables being read
	if !config.IsComponentVersionConfigurationsValid() {
		return Configurations{}, errors.New("mandatory component version of either aws-node, coredns, kube-proxy or cluster-autoscaler not set in config file")
	}

	if !config.IsClusterListConfigurationValid() {
		return Configurations{}, errors.New("one of the clusterlist elements has either ClusterName, AwsRegion, AwsAccount, AwsNodeObject, ClusterAutoscalerObject, KubeProxyObject, CoreDnsObject is missing")
	}

	return config, nil
}

func FileMetadata() (fileName, filePath, fileType string) {
	return FileName, FileType, FilePath
}
