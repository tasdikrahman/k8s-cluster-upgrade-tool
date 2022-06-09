package config

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func TestConfigurations_IsK8sObjectAttributeValid(t *testing.T) {
	tests := []struct {
		name          string
		configuration Configurations
		result        bool
	}{
		{
			name: "when the config passed has all keys and values for awsnode, coredns, clusterautoscaler, kubeproxy present",
			configuration: Configurations{
				ClusterList: []ClusterListConfiguration{
					{
						Name:       "cluster1",
						AwsRegion:  "region",
						AwsAccount: "account",
						AwsNodeObject: K8sObject{
							Name: "aws-node",
							Type: "daemonset",
						},
						ClusterAutoscalerObject: K8sObject{
							Name: "cluster-autoscaler",
							Type: "deployment",
						},
						KubeProxyObject: K8sObject{
							Name: "kube-proxy",
							Type: "daemonset",
						},
						CoreDnsObject: K8sObject{
							Name: "coredns",
							Type: "deployment",
						},
					},
					{
						Name:       "cluster2",
						AwsRegion:  "region",
						AwsAccount: "account",
						AwsNodeObject: K8sObject{
							Name: "aws-node",
							Type: "daemonset",
						},
						ClusterAutoscalerObject: K8sObject{
							Name: "cluster-autoscaler",
							Type: "deployment",
						},
						KubeProxyObject: K8sObject{
							Name: "kube-proxy",
							Type: "daemonset",
						},
						CoreDnsObject: K8sObject{
							Name: "coredns",
							Type: "deployment",
						},
					},
				},
			},
			result: true,
		},
		{
			name: "when the config passed has all keys and values for awsnode, coredns, clusterautoscaler, kubeproxy present but one of the values for the keys is an empty string",
			configuration: Configurations{
				ClusterList: []ClusterListConfiguration{
					{
						Name:       "cluster1",
						AwsRegion:  "region",
						AwsAccount: "account",
						AwsNodeObject: K8sObject{
							Name: "aws-node",
							Type: "daemonset",
						},
						ClusterAutoscalerObject: K8sObject{
							Name: "cluster-autoscaler",
							Type: "deployment",
						},
						KubeProxyObject: K8sObject{
							Name: "kube-proxy",
							Type: "daemonset",
						},
						CoreDnsObject: K8sObject{
							Name: "coredns",
							Type: "deployment",
						},
					},
					{
						Name:       "cluster2",
						AwsRegion:  "region",
						AwsAccount: "account",
						AwsNodeObject: K8sObject{
							Name: "aws-node",
							Type: "daemonset",
						},
						ClusterAutoscalerObject: K8sObject{
							Name: "cluster-autoscaler",
							Type: "deployment",
						},
						KubeProxyObject: K8sObject{
							Name: "kube-proxy",
							Type: "daemonset",
						},
						CoreDnsObject: K8sObject{
							Name: "coredns",
							Type: "",
						},
					},
				},
			},
			result: false,
		},
		{
			name: "when the config passed has one of the k8sObjects keys missing",
			configuration: Configurations{
				ClusterList: []ClusterListConfiguration{
					{
						Name:       "cluster1",
						AwsRegion:  "region",
						AwsAccount: "account",
						AwsNodeObject: K8sObject{
							Name: "aws-node",
							Type: "daemonset",
						},
						ClusterAutoscalerObject: K8sObject{
							Name: "cluster-autoscaler",
							Type: "deployment",
						},
						KubeProxyObject: K8sObject{
							Name: "kube-proxy",
							Type: "daemonset",
						},
						CoreDnsObject: K8sObject{
							Name: "coredns",
							Type: "deployment",
						},
					},
					{
						Name:       "cluster2",
						AwsRegion:  "region",
						AwsAccount: "account",
						AwsNodeObject: K8sObject{
							Name: "aws-node",
							Type: "daemonset",
						},
						ClusterAutoscalerObject: K8sObject{
							Name: "cluster-autoscaler",
							Type: "deployment",
						},
						KubeProxyObject: K8sObject{
							Name: "kube-proxy",
							Type: "daemonset",
						},
					},
				},
			},
			result: false,
		},
		{
			name: "when the config passed has Name attribute value missing",
			configuration: Configurations{
				ClusterList: []ClusterListConfiguration{
					{
						AwsRegion:  "region",
						AwsAccount: "account",
						AwsNodeObject: K8sObject{
							Name: "aws-node",
							Type: "daemonset",
						},
						ClusterAutoscalerObject: K8sObject{
							Name: "cluster-autoscaler",
							Type: "deployment",
						},
						KubeProxyObject: K8sObject{
							Name: "kube-proxy",
							Type: "daemonset",
						},
						CoreDnsObject: K8sObject{
							Name: "coredns",
							Type: "deployment",
						},
					},
					{
						Name:       "cluster2",
						AwsRegion:  "region",
						AwsAccount: "account",
						AwsNodeObject: K8sObject{
							Name: "aws-node",
							Type: "daemonset",
						},
						ClusterAutoscalerObject: K8sObject{
							Name: "cluster-autoscaler",
							Type: "deployment",
						},
						KubeProxyObject: K8sObject{
							Name: "kube-proxy",
							Type: "daemonset",
						},
					},
				},
			},
			result: false,
		},
		{
			name: "when the config passed has AwsRegion attribute value missing",
			configuration: Configurations{
				ClusterList: []ClusterListConfiguration{
					{
						Name:       "cluster-1",
						AwsAccount: "account",
						AwsNodeObject: K8sObject{
							Name: "aws-node",
							Type: "daemonset",
						},
						ClusterAutoscalerObject: K8sObject{
							Name: "cluster-autoscaler",
							Type: "deployment",
						},
						KubeProxyObject: K8sObject{
							Name: "kube-proxy",
							Type: "daemonset",
						},
						CoreDnsObject: K8sObject{
							Name: "coredns",
							Type: "deployment",
						},
					},
					{
						Name:       "cluster2",
						AwsRegion:  "region",
						AwsAccount: "account",
						AwsNodeObject: K8sObject{
							Name: "aws-node",
							Type: "daemonset",
						},
						ClusterAutoscalerObject: K8sObject{
							Name: "cluster-autoscaler",
							Type: "deployment",
						},
						KubeProxyObject: K8sObject{
							Name: "kube-proxy",
							Type: "daemonset",
						},
					},
				},
			},
			result: false,
		},
		{
			name: "when the config passed has AwsAccount attribute value missing",
			configuration: Configurations{
				ClusterList: []ClusterListConfiguration{
					{
						Name:      "cluster-1",
						AwsRegion: "region",
						AwsNodeObject: K8sObject{
							Name: "aws-node",
							Type: "daemonset",
						},
						ClusterAutoscalerObject: K8sObject{
							Name: "cluster-autoscaler",
							Type: "deployment",
						},
						KubeProxyObject: K8sObject{
							Name: "kube-proxy",
							Type: "daemonset",
						},
						CoreDnsObject: K8sObject{
							Name: "coredns",
							Type: "deployment",
						},
					},
					{
						Name:       "cluster2",
						AwsRegion:  "region",
						AwsAccount: "account",
						AwsNodeObject: K8sObject{
							Name: "aws-node",
							Type: "daemonset",
						},
						ClusterAutoscalerObject: K8sObject{
							Name: "cluster-autoscaler",
							Type: "deployment",
						},
						KubeProxyObject: K8sObject{
							Name: "kube-proxy",
							Type: "daemonset",
						},
					},
				},
			},
			result: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.configuration.IsClusterListConfigurationValid(), tt.result)
		})
	}
}

func TestConfigurations_IsClusterNameValid(t *testing.T) {
	tests := []struct {
		name          string
		arg           string
		configuration Configurations
		result        bool
	}{
		{
			name: "when cluster name passed is present in the configuration",
			arg:  "cluster1",
			configuration: Configurations{
				ClusterList: []ClusterListConfiguration{
					{
						Name:       "cluster1",
						AwsRegion:  "region",
						AwsAccount: "account",
					},
					{
						Name:       "cluster2",
						AwsRegion:  "region",
						AwsAccount: "account",
					},
				}},
			result: true,
		},
		{
			name: "when cluster name passed is not present in the configuration",
			arg:  "incorrect-cluster",
			configuration: Configurations{
				ClusterList: []ClusterListConfiguration{
					{
						Name:       "cluster1",
						AwsRegion:  "region",
						AwsAccount: "account",
					},
					{
						Name:       "cluster2",
						AwsRegion:  "region",
						AwsAccount: "account",
					},
				}},
			result: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.result, tt.configuration.IsClusterNameValid(tt.arg))
		})
	}
}

func TestConfigurations_GetAwsAccountAndRegionForCluster(t *testing.T) {
	tests := []struct {
		name             string
		config           Configurations
		arg              string
		awsAccountResult string
		awsRegionResult  string
		err              error
	}{
		{
			name: "returns back the aws account and aws region for the passed cluster when it's found",
			arg:  "cluster1",
			config: Configurations{
				ClusterList: []ClusterListConfiguration{
					{
						Name:       "cluster1",
						AwsRegion:  "region1",
						AwsAccount: "account1",
					},
					{
						Name:       "cluster2",
						AwsRegion:  "region2",
						AwsAccount: "account2",
					},
				}},
			awsAccountResult: "account1",
			awsRegionResult:  "region1",
			err:              nil,
		},
		{
			name: "returns empty strings for aws account and region for the passed cluster name along with an error object",
			arg:  "incorrect-cluster",
			config: Configurations{
				ClusterList: []ClusterListConfiguration{
					{
						Name:       "cluster1",
						AwsRegion:  "region1",
						AwsAccount: "account1",
					},
					{
						Name:       "cluster2",
						AwsRegion:  "region2",
						AwsAccount: "account2",
					},
				}},
			awsAccountResult: "",
			awsRegionResult:  "",
			err:              errors.New("no awsAccount and awsRegion was found for the passed clusterName"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			account, region, err := tt.config.GetAwsAccountAndRegionForCluster(tt.arg)

			assert.Equal(t, account, tt.awsAccountResult)
			assert.Equal(t, region, tt.awsRegionResult)
			assert.Equal(t, err, tt.err)
		})
	}
}

func TestConfigurations_ValidatePassedComponentVersions(t *testing.T) {
	type testArgs struct {
		componentName    string
		componentVersion string
	}
	tests := []struct {
		name   string
		config Configurations
		args   testArgs
		err    error
	}{
		{"when passed component version name is valid and the version to be set matches the config file",
			Configurations{
				Components: ComponentVersionConfigurations{
					CoreDns: "rightvalue", ClusterAutoscaler: "rightvalue", KubeProxy: "rightvalue", AwsNode: "rightvalue"}},
			testArgs{componentName: "coredns", componentVersion: "rightvalue"},
			nil,
		},
		{"when passed component version name is valid and the version to be set doesn't match the config file",
			Configurations{
				Components: ComponentVersionConfigurations{
					CoreDns: "rightvalue", ClusterAutoscaler: "rightvalue", KubeProxy: "rightvalue", AwsNode: "rightvalue"}},
			testArgs{componentName: "coredns", componentVersion: "wrongvalue"},
			errors.New("coredns component version passed doesn't match the version in config, please check the value in config file"),
		},
		{"when passed component version is not valid",
			Configurations{
				Components: ComponentVersionConfigurations{
					CoreDns: "rightvalue", ClusterAutoscaler: "rightvalue", KubeProxy: "rightvalue", AwsNode: "rightvalue"}},
			testArgs{componentName: "foo", componentVersion: "wrongvalue"},
			errors.New("please pass a valid component name from this list [coredns, cluster-autoscaler, kube-proxy, aws-node]"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.ValidatePassedComponentVersions(tt.args.componentName, tt.args.componentVersion)

			assert.Equal(t, err, tt.err)
		})
	}
}

func TestConfigurations_IsComponentVersionConfigurationsValid(t *testing.T) {
	tests := []struct {
		name          string
		configuration Configurations
		result        bool
	}{
		{
			name: "when all the passed component version configurations are present",
			configuration: Configurations{
				Components: ComponentVersionConfigurations{
					CoreDns:           "core-dns-version",
					AwsNode:           "aws-node-version",
					ClusterAutoscaler: "cluster-autoscaler-version",
					KubeProxy:         "kube-proxy-version",
				},
			},
			result: true,
		},
		{
			name: "when one of the required component keys are not passed",
			configuration: Configurations{
				Components: ComponentVersionConfigurations{
					AwsNode:           "aws-node-version",
					ClusterAutoscaler: "cluster-autoscaler-version",
					KubeProxy:         "kube-proxy-version",
				},
			},
			result: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.configuration.IsComponentVersionConfigurationsValid(), tt.result)
		})
	}
}

func TestConfigurations_GetK8sObjectNameAndObjectTypeForCluster(t *testing.T) {
	tests := []struct {
		name                                   string
		configuration                          Configurations
		clusterNameArg, k8sObjectArg           string
		expectedResultName, expectedResultType string
		expectedErr                            error
	}{
		{
			name: "when the cluster name is present and the k8sobject passed is valid",
			configuration: Configurations{
				ClusterList: []ClusterListConfiguration{
					{
						Name:       "cluster1",
						AwsRegion:  "region",
						AwsAccount: "account",
						AwsNodeObject: K8sObject{
							Name: "aws-node",
							Type: "daemonset",
						},
						ClusterAutoscalerObject: K8sObject{
							Name: "cluster-autoscaler",
							Type: "deployment",
						},
						KubeProxyObject: K8sObject{
							Name: "kube-proxy",
							Type: "daemonset",
						},
						CoreDnsObject: K8sObject{
							Name: "coredns",
							Type: "deployment",
						},
					},
					{
						Name:       "cluster2",
						AwsRegion:  "region",
						AwsAccount: "account",
						AwsNodeObject: K8sObject{
							Name: "aws-node",
							Type: "daemonset",
						},
						ClusterAutoscalerObject: K8sObject{
							Name: "cluster-autoscaler",
							Type: "deployment",
						},
						KubeProxyObject: K8sObject{
							Name: "kube-proxy",
							Type: "daemonset",
						},
						CoreDnsObject: K8sObject{
							Name: "coredns",
							Type: "deployment",
						},
					},
				},
			},
			clusterNameArg:     "cluster1",
			k8sObjectArg:       "aws-node",
			expectedResultName: "aws-node",
			expectedResultType: "daemonset",
			expectedErr:        nil,
		},
		{
			name: "when the cluster name is present and the k8sobject passed is invalid",
			configuration: Configurations{
				ClusterList: []ClusterListConfiguration{
					{
						Name:       "cluster1",
						AwsRegion:  "region",
						AwsAccount: "account",
						AwsNodeObject: K8sObject{
							Name: "aws-node",
							Type: "daemonset",
						},
						ClusterAutoscalerObject: K8sObject{
							Name: "cluster-autoscaler",
							Type: "deployment",
						},
						KubeProxyObject: K8sObject{
							Name: "kube-proxy",
							Type: "daemonset",
						},
						CoreDnsObject: K8sObject{
							Name: "coredns",
							Type: "deployment",
						},
					},
					{
						Name:       "cluster2",
						AwsRegion:  "region",
						AwsAccount: "account",
						AwsNodeObject: K8sObject{
							Name: "aws-node",
							Type: "daemonset",
						},
						ClusterAutoscalerObject: K8sObject{
							Name: "cluster-autoscaler",
							Type: "deployment",
						},
						KubeProxyObject: K8sObject{
							Name: "kube-proxy",
							Type: "daemonset",
						},
						CoreDnsObject: K8sObject{
							Name: "coredns",
							Type: "deployment",
						},
					},
				},
			},
			clusterNameArg:     "cluster1",
			k8sObjectArg:       "invalid-arg",
			expectedResultName: "",
			expectedResultType: "",
			expectedErr:        errors.New("please pass any of the components between aws-node, coredns, cluster-autoscaler, kube-proxy"),
		},
		{
			name: "when the cluster name is not present",
			configuration: Configurations{
				ClusterList: []ClusterListConfiguration{
					{
						Name:       "cluster1",
						AwsRegion:  "region",
						AwsAccount: "account",
						AwsNodeObject: K8sObject{
							Name: "aws-node",
							Type: "daemonset",
						},
						ClusterAutoscalerObject: K8sObject{
							Name: "cluster-autoscaler",
							Type: "deployment",
						},
						KubeProxyObject: K8sObject{
							Name: "kube-proxy",
							Type: "daemonset",
						},
						CoreDnsObject: K8sObject{
							Name: "coredns",
							Type: "deployment",
						},
					},
					{
						Name:       "cluster2",
						AwsRegion:  "region",
						AwsAccount: "account",
						AwsNodeObject: K8sObject{
							Name: "aws-node",
							Type: "daemonset",
						},
						ClusterAutoscalerObject: K8sObject{
							Name: "cluster-autoscaler",
							Type: "deployment",
						},
						KubeProxyObject: K8sObject{
							Name: "kube-proxy",
							Type: "daemonset",
						},
						CoreDnsObject: K8sObject{
							Name: "coredns",
							Type: "deployment",
						},
					},
				},
			},
			clusterNameArg:     "invalid cluster",
			k8sObjectArg:       "aws-node",
			expectedResultName: "",
			expectedResultType: "",
			expectedErr:        errors.New("please check if you passed a valid cluster name"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualResultName, actualResultType, actualError := tt.configuration.GetK8sObjectNameAndObjectTypeForCluster(tt.clusterNameArg, tt.k8sObjectArg)

			assert.Equal(t, tt.expectedResultName, actualResultName)
			assert.Equal(t, tt.expectedResultType, actualResultType)
			assert.Equal(t, tt.expectedErr, actualError)
		})
	}
}

func TestRead(t *testing.T) {
	type File struct {
		fileName  string
		fileType  string
		dirName   string
		data      string
		writeFile bool
	}
	tests := []struct {
		name string
		file File
		err  error
	}{
		{"when the config file is present with all the config keys and read successfully",
			File{fileName: "config", fileType: "yaml", dirName: "/tmp", data: "---\ncomponents:\n  aws-node: \"aws-node-version\"\n  cluster-autoscaler: \"cluster-autoscaler-version\"\n  coredns: \"core-dns-version\"\n  kube-proxy: \"kube-proxy-version\"\nclusterlist:\n- Name: \"cluster1\"\n  AwsRegion: \"region1\"\n  AwsAccount: \"account1\"\n  AwsNodeObject:\n    type: \"daemonset\"\n    name: \"aws-node\"\n  ClusterAutoscalerObject:\n    type: \"deployment\"\n    name: \"cluster-autoscaler\"\n  CoreDnsObject:\n    type: \"deployment\"\n    name: \"coredns\"\n  KubeProxyObject:\n    type: \"daemonset\"\n    name: \"kube-proxy\"\n- Name: \"cluster2\"\n  AwsRegion: \"region1\"\n  AwsAccount: \"account1\"\n  AwsNodeObject:\n    type: \"daemonset\"\n    name: \"aws-node\"\n  ClusterAutoscalerObject:\n    type: \"deployment\"\n    name: \"cluster-autoscaler\"\n  CoreDnsObject:\n    type: \"deployment\"\n    name: \"coredns\"\n  KubeProxyObject:\n    type: \"daemonset\"\n    name: \"kube-proxy\"", writeFile: true},
			nil,
		},
		{"when the config file is present and read successfully, but one of the keys for cluster list config is not present with the value",
			File{fileName: "config", fileType: "yaml", dirName: "/tmp", data: "---\ncomponents:\n  aws-node: \"aws-node-version\"\n  cluster-autoscaler: \"cluster-autoscaler-version\"\n  coredns: \"core-dns-version\"\n  kube-proxy: \"kube-proxy-version\"\nclusterlist:\n- Name: \"cluster1\"\n  AwsRegion: \"region1\"\n  AwsAccount: \"account1\"\n- Name: \"cluster2\"\n  AwsRegion: \"region1\"\n  AwsAccount: \"\"\n", writeFile: true},
			errors.New("one of the clusterlist elements has either Name, AwsRegion, AwsAccount, AwsNodeObject, ClusterAutoscalerObject, KubeProxyObject, CoreDnsObject is missing"),
		},
		{"when the config file is present and read successfully, but one of the keys for cluster list config is not present with the key itself",
			File{fileName: "config", fileType: "yaml", dirName: "/tmp", data: "---\ncomponents:\n  aws-node: \"aws-node-version\"\n  cluster-autoscaler: \"cluster-autoscaler-version\"\n  coredns: \"core-dns-version\"\n  kube-proxy: \"kube-proxy-version\"\nclusterlist:\n- Name: \"cluster1\"\n  AwsRegion: \"region1\"\n  AwsAccount: \"account1\"\n- Name: \"cluster2\"\n  AwsRegion: \"region1\"\n", writeFile: true},
			errors.New("one of the clusterlist elements has either Name, AwsRegion, AwsAccount, AwsNodeObject, ClusterAutoscalerObject, KubeProxyObject, CoreDnsObject is missing"),
		},
		{"when the config file is present and read successfully, but kube-proxy config is not present",
			File{fileName: "config", fileType: "yaml", dirName: "/tmp", data: "---\ncomponents:\n  aws-node: \"aws-node-version\"\n  cluster-autoscaler: \"cluster-autoscaler-version\"\n  coredns: \"core-dns-version\"\nclusterlist:\n- Name: \"cluster1\"\n  AwsRegion: \"region1\"\n  AwsAccount: \"account1\"\n- Name: \"cluster2\"\n  AwsRegion: \"region1\"\n  AwsAccount: \"account1\"\n", writeFile: true},
			errors.New("mandatory component version of either aws-node, coredns, kube-proxy or cluster-autoscaler not set in config file"),
		},
		{"when the config file is present and read successfully, but aws-node config is not present",
			File{fileName: "config", fileType: "yaml", dirName: "/tmp", data: "---\ncomponents:\n  kube-proxy: \"kube-proxy-version\"\n  cluster-autoscaler: \"cluster-autoscaler-version\"\n  coredns: \"core-dns-version\"\nclusterlist:\n- Name: \"cluster1\"\n  AwsRegion: \"region1\"\n  AwsAccount: \"account1\"\n- Name: \"cluster2\"\n  AwsRegion: \"region1\"\n  AwsAccount: \"account1\"\n", writeFile: true},
			errors.New("mandatory component version of either aws-node, coredns, kube-proxy or cluster-autoscaler not set in config file"),
		},
		{"when the config file is present and read successfully, but cluster-autoscaler config is not present",
			File{fileName: "config", fileType: "yaml", dirName: "/tmp", data: "---\ncomponents:\n  kube-proxy: \"kube-proxy-version\"\n  aws-node: \"aws-node-version\"\n  coredns: \"core-dns-version\"\nclusterlist:\n- Name: \"cluster1\"\n  AwsRegion: \"region1\"\n  AwsAccount: \"account1\"\n- Name: \"cluster2\"\n  AwsRegion: \"region1\"\n  AwsAccount: \"account1\"\n", writeFile: true},
			errors.New("mandatory component version of either aws-node, coredns, kube-proxy or cluster-autoscaler not set in config file"),
		},
		{"when the config file is present and read successfully, but coredns config is not present",
			File{fileName: "config", fileType: "yaml", dirName: "/tmp", data: "---\ncomponents:\n  kube-proxy: \"kube-proxy-version\"\n  aws-node: \"aws-node-version\"\n  cluster-autoscaler: \"cluster-autoscaler-version\"\nclusterlist:\n- Name: \"cluster1\"\n  AwsRegion: \"region1\"\n  AwsAccount: \"account1\"\n- Name: \"cluster2\"\n  AwsRegion: \"region1\"\n  AwsAccount: \"account1\"\n", writeFile: true},
			errors.New("mandatory component version of either aws-node, coredns, kube-proxy or cluster-autoscaler not set in config file"),
		},
		{"when the config file is not present",
			File{fileName: "config", fileType: "yaml", dirName: "/tmp", data: "", writeFile: false},
			errors.New("error finding config file. Does it exist? Please create it in $HOME/.k8s-cluster-upgrade-tool/config.yaml if not"),
		},
		{"when the config file is present, but reading fails as data type inside is not yaml",
			File{fileName: "config", fileType: "yaml", dirName: "/tmp", data: "foo baz", writeFile: true},
			errors.New("error reading from config file"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.file.writeFile {
				fileContent := []byte(fmt.Sprintf(tt.file.data))
				err := ioutil.WriteFile(fmt.Sprintf("%s/%s.%s", tt.file.dirName, tt.file.fileName, tt.file.fileType),
					fileContent, 0644)
				if err != nil {
					log.Fatal("error writing to temp config file for running tests")
				}
				defer os.Remove(fmt.Sprintf("%s/%s.%s", tt.file.dirName, tt.file.fileName, tt.file.fileType))
			}

			_, err := Read(tt.file.fileName, tt.file.fileType, tt.file.dirName)

			assert.Equal(t, err, tt.err)
		})
	}
}

func TestFileMetadata(t *testing.T) {
	t.Run("returns the correct path, filetype and directory", func(t *testing.T) {
		gotFileName, gotFileType, gotFilePath := FileMetadata()

		assert.Equal(t, gotFileName, "config")
		assert.Equal(t, gotFileType, "yaml")
		assert.Equal(t, gotFilePath, "$HOME/.k8s-cluster-upgrade-tool")
	})
}
