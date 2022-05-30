package aws

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	"log"
)

// TODO: Improve the modelling of cluster and awsinstances to be in the appropriate packages.
type AutoScalingGroup struct {
	Instances        AwsInstances
	DesiredInstances int
	MinInstances     int
	MaxInstances     int
	AsgName          string
}

type UpdateAutoscalingGroupInterface interface {
	UpdateAutoScalingGroupCount(ctx context.Context, cfg aws.Config) (*autoscaling.UpdateAutoScalingGroupOutput, error)
}

type AutoScalingGroupClient struct {
	Asg AutoScalingGroup
}

func (a *AutoScalingGroupClient) UpdateAutoScalingGroupCount(ctx context.Context, cfg aws.Config) (*autoscaling.UpdateAutoScalingGroupOutput, error) {
	autoscalingAwsClient := autoscaling.NewFromConfig(cfg)
	input := &autoscaling.UpdateAutoScalingGroupInput{
		AutoScalingGroupName: aws.String(a.Asg.AsgName),
		MaxSize:              aws.Int32(int32(a.Asg.DesiredInstances)),
	}

	result, err := autoscalingAwsClient.UpdateAutoScalingGroup(ctx, input)
	if err != nil {
		log.Println("There was an error updating the autoscaling group's max instance count")
		return &autoscaling.UpdateAutoScalingGroupOutput{}, err
	}
	return result, nil
}

type AutoscalingGroupUpdater struct {
	UpdateAutoscalingGroupInterface
}

// Update method updates the current ASG count for max to the current desired count to prevent the ASG being drained and tainted
// to prevent from getting autoscaled during the upgrade process
// TODO make this general purpose instead of just updating the Max size to current desired size which it is doing right now
func (a *AutoscalingGroupUpdater) Update(ctx context.Context, awsConfig aws.Config) (*autoscaling.UpdateAutoScalingGroupOutput, error) {
	result, err := a.UpdateAutoScalingGroupCount(ctx, awsConfig)
	if err != nil {
		return nil, err
	}
	return result, nil
}
