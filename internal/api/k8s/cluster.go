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
	"k8s.io/client-go/util/retry"
	"log"
	"os/exec"
	"path/filepath"
	"strings"
)

// ParseComponentImage takes in the full container image and returns back the container tag or the container image
// based on the argument received
func ParseComponentImage(kubectlExecOutput string, imageSection string) (string, error) {
	if imageSection == "imageTag" {
		return strings.Split(kubectlExecOutput, ":")[1], nil
	} else if imageSection == "imagePrefix" {
		return strings.Split(kubectlExecOutput, ":")[0], nil
	} else {
		return "", errors.New("invalid imageSection Passed")
	}
}

// TODO replace this with client-go
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

// TODO replace this with client-go
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

// TODO replace this with client-go
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

// KubeClientInit returns back clientSet
func KubeClientInit(kubeContext string) (*kubernetes.Clientset, error) {
	var kubeConfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeConfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeConfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	config, err := buildConfigFromFlags(kubeContext, *kubeConfig)
	if err != nil {
		return &kubernetes.Clientset{}, errors.New("error building the config for building the client-set for client-go")
	}

	// create the clientset
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return &kubernetes.Clientset{}, errors.New("error building the client-set for client-go")
	}
	return clientSet, nil
}

// GetContainerImageForK8sObject is used to return  the container image from for the object
// Supports deployment and Daemonsets as of now for apps/v1 api
// The clienset would have already been initialized with the specific k8s context to be used with
//
// Usage:
// kubeClient, _ := k8s.KubeClientInit("cluster-name")
// containerImage, _ := k8s.GetContainerImageForK8sObject(kubeClient, "aws-node", "daemonset", "kube-system")
func GetContainerImageForK8sObject(k8sClient kubernetes.Interface, k8sObjectName, k8sObject, namespace string) (string, error) {
	switch k8sObject {
	case "deployment":
		// NOTE: Not targeting other api versions for the objects as of now.
		deployment, err := k8sClient.AppsV1().Deployments(namespace).Get(context.TODO(), k8sObjectName, metav1.GetOptions{})
		if k8sErrors.IsNotFound(err) {
			return "", fmt.Errorf("Deployment %s in namespace %s not found\n", k8sObjectName, namespace)
		} else if statusError, isStatus := err.(*k8sErrors.StatusError); isStatus {
			return "", fmt.Errorf("Error getting deployment %s in namespace %s: %v\n",
				k8sObjectName, namespace, statusError.ErrStatus.Message)
		} else if err != nil {
			return "", fmt.Errorf("there was an error while retrieving the container image")
		}

		// NOTE: This assumes there is only one container in the k8s object, which is true for the components for us at moment
		return deployment.Spec.Template.Spec.Containers[0].Image, nil
	case "daemonset":
		// NOTE: Not targeting other api versions for the objects as of now.
		daemonSet, err := k8sClient.AppsV1().DaemonSets(namespace).Get(context.TODO(), k8sObjectName, metav1.GetOptions{})
		if k8sErrors.IsNotFound(err) {
			return "", fmt.Errorf("daemonset %s in namespace %s not found\n", k8sObjectName, namespace)
		} else if statusError, isStatus := err.(*k8sErrors.StatusError); isStatus {
			return "", fmt.Errorf(fmt.Sprintf("Error getting daemonset %s in namespace %s: %v\n",
				k8sObjectName, namespace, statusError.ErrStatus.Message))
		} else if err != nil {
			return "", fmt.Errorf("there was an error while retrieving the container image")
		}

		// NOTE: This assumes there is only one container in the k8s object, which is true for the components for us at moment
		return daemonSet.Spec.Template.Spec.Containers[0].Image, nil
	default:
		return "", fmt.Errorf("please choose between Daemonset or Deployment k8sobject as they are currently supported")
	}
}

// SetK8sObjectImage will set the image version for the deployment/daemonset object requested to update for
func SetK8sObjectImage(k8sClient kubernetes.Interface, k8sObject, k8sObjectName, containerName, containerImage, k8sNamespace string) error {
	containerFound := false
	switch k8sObject {
	case "deployment":
		retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
			// Retrieve the latest version of Deployment before attempting update
			// RetryOnConflict uses exponential backoff to avoid exhausting the apiserver
			result, getErr := k8sClient.AppsV1().Deployments(k8sNamespace).Get(context.TODO(), k8sObjectName, metav1.GetOptions{})
			if getErr != nil {
				return fmt.Errorf("failed to get latest version of Deployment: %v", getErr)
			}

			for containerNumber, candidateContainer := range result.Spec.Template.Spec.Containers {
				if candidateContainer.Name == containerName {
					containerFound = true
					result.Spec.Template.Spec.Containers[containerNumber].Image = containerImage // update container image
				}
			}
			if !containerFound {
				return fmt.Errorf("container %s was not found in the deployment object, skipping update", containerName)
			}

			_, updateErr := k8sClient.AppsV1().Deployments(k8sNamespace).Update(context.TODO(), result, metav1.UpdateOptions{})
			return updateErr
		})
		if retryErr != nil {
			return fmt.Errorf("container image update failed: %v", retryErr)
		}
		return nil
	case "daemonset":
		retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
			// Retrieve the latest version of Deployment before attempting update
			// RetryOnConflict uses exponential backoff to avoid exhausting the apiserver
			result, getErr := k8sClient.AppsV1().DaemonSets(k8sNamespace).Get(context.TODO(), k8sObjectName, metav1.GetOptions{})
			if getErr != nil {
				return fmt.Errorf("failed to get latest version of Daemonset: %v", getErr)
			}

			for containerNumber, candidateContainer := range result.Spec.Template.Spec.Containers {
				if candidateContainer.Name == containerName {
					containerFound = true
					result.Spec.Template.Spec.Containers[containerNumber].Image = containerImage // update container image
				}
			}
			if !containerFound {
				return fmt.Errorf("container %s was not found in the daemonset object, skipping update", containerName)
			}
			_, updateErr := k8sClient.AppsV1().DaemonSets(k8sNamespace).Update(context.TODO(), result, metav1.UpdateOptions{})
			return updateErr
		})
		if retryErr != nil {
			return fmt.Errorf("container image update failed: %v", retryErr)
		}
		return nil
	default:
		return errors.New("please pass the k8sObject to be from daemonset or deployment")
	}
}
