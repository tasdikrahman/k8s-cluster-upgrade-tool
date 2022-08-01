package aws

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-sdk-go-v2/aws"
	"log"
	"os/exec"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go/aws/awserr"

	"github.com/deliveryhero/k8s-cluster-upgrade-tool/internal/api/k8s"
)

// TODO: Make the domain modelling cleaner and compose the types together which are created alongside one another
type AwsInstance struct {
	InstanceId string
	PrivateDNS string
	AsgName    string
}

type AwsInstances []AwsInstance

// TODO Add a spec for this
func (a *AwsInstances) AppendInstance(o AwsInstance) {
	*a = append(*a, o)
}

// Count returns the number of AWS instances currently running in the ASG, a quick way to gather the desired count from
// instance map created already
func (a AwsInstances) Count() int {
	return len(a)
}

// TODO Add a spec for this
func (a AwsInstances) PrettyPrint() {
	for _, instance := range a {
		jsonData, err := json.Marshal(&instance)
		if err != nil {
			log.Println("Error with marshaling data while printing instance")
			log.Fatal(err)
		}
		log.Println(string(jsonData))
	}
}

// GetInstancesForASG is a helper function, which interacts with the AWS SDK taking the input of the asgname, awsregion
// and awsprofile and then calls the DescribeAutoScalingGroups API to get the instances of the ASG and then calling the
// DescribeInstances API to map the private DNS of the instances and then store it, which will then be used to feed
// kubectl drain command
// TODO have this method use the generic config getter helper created in internal/api/aws package to reduce duplication
// TODO break this method into two helper methods
// TODO add specs for this
func (a *AwsInstances) GetInstancesForASG(cfg aws.Config, asgName string, awsRegion string, awsProfile string) {
	autoscalingAwsClient := autoscaling.NewFromConfig(cfg)
	input := &autoscaling.DescribeAutoScalingGroupsInput{
		AutoScalingGroupNames: []string{
			asgName,
		},
	}

	// TODO Add non happy path to give a clear error message if the ASG passed is not present in the AWS PROFILE passed
	// TODO Describe the tags of the autoscaling group to check the cluster tag and match it and exit out if they don't match
	describeAutoScalingGroupsResult, err := autoscalingAwsClient.DescribeAutoScalingGroups(context.TODO(), input)
	if err != nil {
		log.Println(err.Error())
		return
	}

	var awsInstanceIds []string
	for _, instance := range describeAutoScalingGroupsResult.AutoScalingGroups[0].Instances {
		awsInstanceIds = append(awsInstanceIds, *instance.InstanceId)
	}

	ec2AwsClient := ec2.NewFromConfig(cfg)
	ec2Input := &ec2.DescribeInstancesInput{
		InstanceIds: awsInstanceIds,
	}

	ec2result, ec2err := ec2AwsClient.DescribeInstances(context.TODO(), ec2Input)
	if ec2err != nil {
		if ec2err, ok := err.(awserr.Error); ok {
			switch ec2err.Code() {
			default:
				log.Println(ec2err.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			log.Println(ec2err.Error())
		}
		return
	}
	for _, reservations := range ec2result.Reservations {
		if reservations.Instances[0].State.Name == "running" {
			awsInstance := AwsInstance{
				InstanceId: *reservations.Instances[0].InstanceId,
				PrivateDNS: *reservations.Instances[0].PrivateDnsName,
				AsgName:    asgName,
			}
			a.AppendInstance(awsInstance)
		} else {
			log.Fatal("One or many of the instances are not in the running state, please check the ASG on console")
		}
	}
}

// TODO Add a spec for this
func (a AwsInstances) TaintNodes() error {
	for _, instance := range a {
		log.Printf("Tainting node: %s\n", instance.PrivateDNS)
		args := strings.Fields(k8s.KubectlTaintNodeCommand(instance.PrivateDNS))

		output, err := exec.Command(args[0], args[1:]...).Output()
		if err != nil {
			log.Fatal("There was an error while tainting the node: ", err)
			return err
		}
		log.Printf("taint output: \n %s", output)
	}
	return nil
}

// TODO Add a spec for this
func (a AwsInstances) DrainNodes() error {
	for _, instance := range a {
		log.Printf("Draining node: %s\n", instance.PrivateDNS)
		args := strings.Fields(k8s.KubectlDrainNodeCommand(instance.PrivateDNS))

		output, err := exec.Command(args[0], args[1:]...).Output()
		if err != nil {
			log.Fatal("There was an error while draining the node: ", err)
			return err
		}
		log.Printf("drain output: \n %s", output)

	}
	return nil
}
