package k8s

import (
	"errors"
	"fmt"
	"log"
	"os/exec"
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
