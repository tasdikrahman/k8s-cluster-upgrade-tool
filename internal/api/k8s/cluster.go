package k8s

import (
	"context"
	"errors"
	"flag"
	"fmt"
	k8sErrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"log"
	"os/exec"
	"path/filepath"
	"strings"
)

func ParseComponentImage(kubectlExecOutput string, imageSection string) (string, error) {
	if imageSection == "imageTag" {
		return strings.Trim(strings.Split(kubectlExecOutput, ":")[1], "'"), nil
	} else if imageSection == "imagePrefix" {
		return strings.Trim(strings.Split(kubectlExecOutput, ":")[0], "'"), nil
	} else {
		return "", errors.New("invalid imageSection Passed")
	}
}

// TODO add spec for this
func KubectlGetImageCommand(k8sObject, component, namespace string) string {
	return fmt.Sprintf(`
	kubectl
	get
	%s
	%s
	--namespace %s
	-o=jsonpath='{$.spec.template.spec.containers[:1].image}'
	`, k8sObject, component, namespace)
}

// TODO add spec for this
func KubectlSetImageCommand(k8sObject, componentName, containerImage, namespace string) string {
	return fmt.Sprintf(`
	kubectl
	set
	image
	%s
	--namespace %s
	%s=%s
	`, k8sObject, componentName, containerImage, namespace)
}

// TODO add spec for this
func KubectlTaintNodeCommand(node string) string {
	// Format: kubectl taint nodes NODE key=value:NoSchedule
	return fmt.Sprintf(`
	kubectl
	taint
	nodes
	%s
	taintkey=k8s-cluster-upgrade-tool:NoSchedule
	`, node)
}

// TODO add spec for this
func KubectlDrainNodeCommand(node string) string {
	// Format: kubectl drain --ignore-daemonsets --force --delete-local-data <node name>
	return fmt.Sprintf(`
	kubectl
	drain
	--ignore-daemonsets
	--force
	--delete-local-data
	%s
	`, node)
}

// TODO add spec for this
func SetK8sContext(clusterName string) {
	command := "kubectl"
	arg01 := "config"
	arg02 := "use-context"

	// TODO: change this to use client-go
	cmd := exec.Command(command, arg01, arg02, clusterName)
	err := cmd.Run()
	if err != nil {
		log.Fatalln("Error setting kube context to the cluster selected")
	}
}

// buildConfigFromFlags returns the config using which the client will be initialized with the k8s context we want to use
func buildConfigFromFlags(context, kubeconfigPath string) (*rest.Config, error) {
	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfigPath},
		&clientcmd.ConfigOverrides{
			CurrentContext: context,
		}).ClientConfig()
}

// GetContainerImageForK8sObject is used to return  the container image from for the object
// Supports deployment and Daemonsets as of now for apps/v1 api
//
// Usage:
// image, err := k8s.GetContainerImageForK8sObject("cluster-name", "aws-node", "daemonset", "kube-system")
func GetContainerImageForK8sObject(kubeContext, k8sObjectName, k8sObject, namespace string) (string, error) {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	config, err := buildConfigFromFlags(kubeContext, *kubeconfig)
	if err != nil {
		return "", errors.New("error building the config for building the client-set for client-go")
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return "", errors.New("error building the client-set for client-go")
	}

	if k8sObject == "deployment" {
		// NOTE: Not targeting other api versions for the objects as of now.
		deployment, err := clientset.AppsV1().Deployments(namespace).Get(context.TODO(), k8sObjectName, metav1.GetOptions{})
		if k8sErrors.IsNotFound(err) {
			return "", errors.New(fmt.Sprintf("Deployment %s in namespace %s not found\n", k8sObjectName, namespace))
		} else if statusError, isStatus := err.(*k8sErrors.StatusError); isStatus {
			return "", errors.New(fmt.Sprintf("Error getting deployment %s in namespace %s: %v\n",
				k8sObjectName, namespace, statusError.ErrStatus.Message))
		} else if err != nil {
			return "", errors.New("there was an error while retrieving the container image")
		}

		// NOTE: This assumes there is only one container in the k8s object, which is true for the components for us at moment
		return deployment.Spec.Template.Spec.Containers[0].Image, nil
	} else if k8sObject == "daemonset" {
		// NOTE: Not targeting other api versions for the objects as of now.
		daemonset, err := clientset.AppsV1().DaemonSets(namespace).Get(context.TODO(), k8sObjectName, metav1.GetOptions{})
		if k8sErrors.IsNotFound(err) {
			return "", errors.New(fmt.Sprintf("Daemonset %s in namespace %s not found\n", k8sObjectName, namespace))
		} else if statusError, isStatus := err.(*k8sErrors.StatusError); isStatus {
			return "", errors.New(fmt.Sprintf("Error getting deployment %s in namespace %s: %v\n",
				k8sObjectName, namespace, statusError.ErrStatus.Message))
		} else if err != nil {
			return "", errors.New("there was an error while retrieving the container image")
		}

		// NOTE: This assumes there is only one container in the k8s object, which is true for the components for us at moment
		return daemonset.Spec.Template.Spec.Containers[0].Image, nil
	} else {
		return "", errors.New("please choose between Daemonset or Deployment k8sobject as they are currently supported")
	}
}
