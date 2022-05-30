package cmd

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/spf13/cobra"
	"log"

	toolConfig "k8s-cluster-upgrade-tool/config"
	"k8s-cluster-upgrade-tool/internal/api/aws"
)

var DryRunFlag bool

var nodeTaintAndDrainCmd = &cobra.Command{
	Use:   "taint-and-drain-asg",
	Short: "Taints and drains nodes from ASG",
	Long: `taint-and-drain-asg helps you taint and drain an ASG in an automated fashion by taking input of the ASG name, nodes of
which you would want to drain and taint later.

It first sets the max instance count of the ASG to the current desired count.
taints the nodes in the ASG
drains the nodes in the ASG

Usage:
$ k8s-cluster-upgrade-tool taint-and-drain-asg -c=CLUSTER_NAME -a=ASG_NAME

Example:
$ k8s-cluster-upgrade-tool taint-and-drain-asg -c=valid-cluster-name -a=valid-cluster-name-spot-hash
$ k8s-cluster-upgrade-tool taint-and-drain-asg -c=valid-cluster-name -a=valid-cluster-name-spot-hash --dry-run=false

For a managed node group, we need to pass the exact ASG resource name, rather than the one which shows up on the EKS console
$ k8s-cluster-upgrade-tool taint-and-drain-asg -c=valid-cluster-name -a=valid-cluster-name-foo-name // incorrect
$ k8s-cluster-upgrade-tool taint-and-drain-asg -c=valid-cluster-name -a=eks-hash-value-asg-name // correct
`,
	Run: func(cmd *cobra.Command, args []string) {
		cluster, _ := cmd.Flags().GetString("cluster")
		asg, _ := cmd.Flags().GetString("autoscaling-group")
		dryRun, _ := cmd.Flags().GetBool("dry-run")

		// validate the cluster name and mapping if it's present
		if toolConfig.Configuration.IsClusterNameValid(cluster) {
			_, _, err := toolConfig.Configuration.GetAwsAccountAndRegionForCluster(cluster)
			if err == nil {
				log.Println("Setting kubernetes context to", cluster)
				setK8sContext(cluster)
			}
		} else {
			log.Fatal("Please pass a valid clusterName or check if the AWS account has a mapping inside the tool for the account and the region")
		}

		// storing all the instances with their private DNS's for the passed ASG for the AWS profile mapped for the cluster passed
		awsAccount, awsRegion, _ := toolConfig.Configuration.GetAwsAccountAndRegionForCluster(cluster)

		// create aws config
		awsGetterObj := &aws.ConfigGetter{ConfigClientInterface: &aws.Config{}}
		cfg, err := awsGetterObj.GetConfig(context.TODO(), config.WithRegion(awsRegion), config.WithSharedConfigProfile(awsAccount))
		if err != nil {
			log.Fatal("there was an error while initializing the aws config, please check your aws credentials")
		}

		awsInstances := aws.AwsInstances{}
		awsInstances.GetInstancesForASG(cfg, asg, awsRegion, awsAccount)

		if dryRun {
			log.Println("Running taint and drain nodes command in dry mode")
			log.Println("Instances which are going to be tainted and drained from the ASG passed")
			awsInstances.PrettyPrint()
		} else {
			log.Println("Running taint and drain command in non-dry mode")

			// add logic Print the instances which are going to be taint and drained
			log.Println("Instances which are going to be tainted and drained from the ASG passed")
			awsInstances.PrettyPrint()

			// add logic which modifies the ASG's Max size to the current desired count to prevent the ASG to scaling up
			asgObject := aws.AutoScalingGroup{
				AsgName:          asg,
				Instances:        awsInstances,
				DesiredInstances: awsInstances.Count(),
			}
			awsAsgClient := &aws.AutoScalingGroupClient{Asg: asgObject}
			// call the autoscaling group update call
			awsUpdateAsgObj := &aws.AutoscalingGroupUpdater{
				UpdateAutoscalingGroupInterface: awsAsgClient,
			}
			_, err := awsUpdateAsgObj.Update(context.TODO(), cfg)
			if err != nil {
				log.Fatal("Updation of the Autoscaling group to make the maximum nodes to be equal to the current number of nodes failed," +
					" skipping, tainting and draining of the ASG")
			}
			log.Printf("The ASG's max size was set to the current desired size, current max size after updation: %d\n",
				awsInstances.Count())

			// iterate over the nodes now to run kubectl taint
			err = awsInstances.TaintNodes()
			if err != nil {
				log.Printf("Error tainting the nodes %s", err)
			}

			// iterate over the nodes now to run kubectl drain
			err = awsInstances.DrainNodes()
			if err != nil {
				log.Printf("Error draining the nodes %s", err)
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(nodeTaintAndDrainCmd)

	nodeTaintAndDrainCmd.Flags().StringP("cluster", "c", "",
		"Example cluster name input valid-cluster-name, check with team for a full list of valid clusters")
	nodeTaintAndDrainCmd.Flags().StringP("autoscaling-group", "a", "",
		"Example cluster name input being valid-cluster-name and the asg name passed being valid-cluster-name-spot-hash")
	nodeTaintAndDrainCmd.Flags().BoolVar(&DryRunFlag, "dry-run", true,
		"will only show the nodes which will be fed to taint and drain")
	//nolint
	nodeTaintAndDrainCmd.MarkFlagRequired("cluster")
	//nolint
	nodeTaintAndDrainCmd.MarkFlagRequired("autoscaling-group")
}
