package k8s

import (
	"errors"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseComponentImage(t *testing.T) {
	type args struct {
		kubectlExecOutput string
		imageSection      string
	}
	tests := []struct {
		name string
		args args
		want string
		err  error
	}{
		{"when getComponentImageTag is passed with a valid output and imageSection and it returns the image tag",
			args{"my-hash.dkr.ecr.eu-west-1.amazonaws.com/amazon-k8s-cni:my-version", "imageTag"},
			"my-version", nil},
		{"when getComponentImageTag is passed with a valid output and imagePrefix and it returns the image tag",
			args{"my-hash.dkr.ecr.eu-west-1.amazonaws.com/amazon-k8s-cni:my-version", "imagePrefix"},
			"my-hash.dkr.ecr.eu-west-1.amazonaws.com/amazon-k8s-cni", nil},
		{"when getComponentImageTag is passed with a valid output and imagePrefix and it returns the image tag",
			args{"my-hash.dkr.ecr.eu-west-1.amazonaws.com/amazon-k8s-cni:my-version", "foo"},
			"", errors.New("invalid imageSection Passed")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseComponentImage(tt.args.kubectlExecOutput, tt.args.imageSection)
			assert.Equal(t, got, tt.want)
			assert.Equal(t, tt.err, err)
		})
	}
}

func TestGetContainerImageForK8sObjectWhenK8sObjectIsDeployment(t *testing.T) {
	type deploymentArgs struct {
		k8sObject     string
		k8sObjectName string
		kubeContext   string
		namespace     string
		deployment    *appsv1.Deployment
	}
	tests := []struct {
		name   string
		args   deploymentArgs
		err    error
		output string
	}{
		{
			name: "When the object is of type deployment, the objectname is cluster-autoscaler, object exists and returns back the image",
			args: deploymentArgs{k8sObject: "deployment", k8sObjectName: "cluster-autoscaler", kubeContext: "test-context", namespace: "kube-system",
				deployment: &appsv1.Deployment{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "cluster-autoscaler",
						Namespace: "kube-system",
					},
					Spec: appsv1.DeploymentSpec{
						Template: corev1.PodTemplateSpec{
							Spec: corev1.PodSpec{
								Containers: []corev1.Container{
									{
										Image: "cluster-autoscaler:v1.0.0",
									},
								},
							},
						},
					},
				}},
			output: "cluster-autoscaler:v1.0.0",
			err:    nil,
		},
		{
			name: "When the object is of type deployment, the objectname is cluster-autoscaler, object doesn't exist, returns back error",
			args: deploymentArgs{k8sObject: "deployment", k8sObjectName: "cluster-autoscaler", kubeContext: "test-context", namespace: "kube-system",
				deployment: &appsv1.Deployment{}},
			output: "",
			err:    errors.New("Deployment cluster-autoscaler in namespace kube-system not found\n"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := fake.NewSimpleClientset(tt.args.deployment)

			got, err := GetContainerImageForK8sObject(client, tt.args.k8sObjectName, tt.args.k8sObject, tt.args.namespace)

			assert.Equal(t, tt.output, got)
			assert.Equal(t, tt.err, err)
		})
	}
}

func TestGetContainerImageForK8sObjectWhenK8sObjectIsDaemonSet(t *testing.T) {
	type daemonSetArgs struct {
		k8sObject     string
		k8sObjectName string
		kubeContext   string
		namespace     string
		daemonSet     *appsv1.DaemonSet
	}
	tests := []struct {
		name   string
		args   daemonSetArgs
		err    error
		output string
	}{
		{
			name: "When the object is of type daemonset, the objectname is aws-node, object exists and returns back the image",
			args: daemonSetArgs{k8sObject: "daemonset", k8sObjectName: "aws-node", kubeContext: "test-context", namespace: "kube-system",
				daemonSet: &appsv1.DaemonSet{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "aws-node",
						Namespace: "kube-system",
					},
					Spec: appsv1.DaemonSetSpec{
						Template: corev1.PodTemplateSpec{
							Spec: corev1.PodSpec{
								Containers: []corev1.Container{
									{
										Image: "aws-node:v1.0.0",
									},
								},
							},
						},
					},
				}},
			output: "aws-node:v1.0.0",
			err:    nil,
		},
		{
			name: "When the object is of type daemonset, the objectname is aws-node, object doesn't exist, returns back error",
			args: daemonSetArgs{k8sObject: "daemonset", k8sObjectName: "aws-node", kubeContext: "test-context", namespace: "kube-system",
				daemonSet: &appsv1.DaemonSet{}},
			output: "",
			err:    errors.New("daemonset aws-node in namespace kube-system not found\n"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := fake.NewSimpleClientset(tt.args.daemonSet)

			got, err := GetContainerImageForK8sObject(client, tt.args.k8sObjectName, tt.args.k8sObject, tt.args.namespace)

			assert.Equal(t, tt.output, got)
			assert.Equal(t, tt.err, err)
		})
	}
}

func TestSetK8sObjectImageWhenObjectIsDeployment(t *testing.T) {
	type deploymentArgs struct {
		k8sObject            string
		k8sObjectName        string
		kubeContext          string
		namespace            string
		targetContainerImage string
		deployment           *appsv1.Deployment
	}
	tests := []struct {
		name string
		args deploymentArgs
		err  error
	}{
		{
			name: "when the deployment is present and the update call goes through",
			args: deploymentArgs{k8sObject: "deployment", k8sObjectName: "cluster-autoscaler", kubeContext: "test-kubecontext",
				namespace: "kube-system", targetContainerImage: "v1.targetversion",
				deployment: &appsv1.Deployment{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "cluster-autoscaler",
						Namespace: "kube-system",
					},
					Spec: appsv1.DeploymentSpec{
						Template: corev1.PodTemplateSpec{
							Spec: corev1.PodSpec{
								Containers: []corev1.Container{
									{
										Image: "cluster-autoscaler:v1.0.0",
									},
								},
							},
						},
					},
				}},
			err: nil,
		},
		{
			name: "when the deployment is not present and the update call fails",
			args: deploymentArgs{k8sObject: "deployment", k8sObjectName: "cluster-autoscaler", kubeContext: "test-kubecontext",
				namespace: "kube-system", targetContainerImage: "v1.targetversion",
				deployment: &appsv1.Deployment{}},
			err: errors.New("container image update failed: failed to get latest version of Deployment: deployments.apps \"cluster-autoscaler\" not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := fake.NewSimpleClientset(tt.args.deployment)
			err := SetK8sObjectImage(client, tt.args.k8sObject, tt.args.k8sObjectName, tt.args.targetContainerImage, tt.args.namespace)

			assert.Equal(t, tt.err, err)
		})
	}
}

func TestSetK8sObjectImageWhenObjectIsDaemonSet(t *testing.T) {
	type daemonSetArgs struct {
		k8sObject            string
		k8sObjectName        string
		kubeContext          string
		namespace            string
		targetContainerImage string
		daemonSet            *appsv1.DaemonSet
	}
	tests := []struct {
		name string
		args daemonSetArgs
		err  error
	}{
		{
			name: "when the daemonset is present and the update call goes through",
			args: daemonSetArgs{k8sObject: "daemonset", k8sObjectName: "aws-node", kubeContext: "test-kubecontext",
				namespace: "kube-system", targetContainerImage: "v1.targetversion",
				daemonSet: &appsv1.DaemonSet{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "aws-node",
						Namespace: "kube-system",
					},
					Spec: appsv1.DaemonSetSpec{
						Template: corev1.PodTemplateSpec{
							Spec: corev1.PodSpec{
								Containers: []corev1.Container{
									{
										Image: "aws-node:v1.0.0",
									},
								},
							},
						},
					},
				}},
			err: nil,
		},
		{
			name: "when the daemonset is not present and the update call fails",
			args: daemonSetArgs{k8sObject: "daemonset", k8sObjectName: "aws-node", kubeContext: "test-kubecontext",
				namespace: "kube-system", targetContainerImage: "v1.targetversion",
				daemonSet: &appsv1.DaemonSet{}},
			err: errors.New("container image update failed: failed to get latest version of Daemonset: daemonsets.apps \"aws-node\" not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := fake.NewSimpleClientset(tt.args.daemonSet)
			err := SetK8sObjectImage(client, tt.args.k8sObject, tt.args.k8sObjectName, tt.args.targetContainerImage, tt.args.namespace)

			assert.Equal(t, tt.err, err)
		})
	}
}
