## k8s-cluster-upgrade-tool

The tool allows you to 
- check for the components installed in the cluster and see whether everything is running in the required version or not and if not
- Set the component running in the cluster to the desired version.
- taint and drain the nodes for an ASG

### Pre-requisite setup

- The cli supports both mac and Linux machines at the moment. You can download the respective binaries from the releases page.
- You have logged into the particular `AWS_PROFILE`, using the authz/authn mechanism, and your user has permissions to modify ASG's for your account and region.
- You are able to run, `use-context` (with the presence of `~/.kube/config` on your machine) for your cluster, say for example you want to interact with valid-cluster-name cluster, you are able to run the below command

```
$ kubectl config use-context valid-cluster-name
Switched to context "valid-cluster-name".
```
- Copy the config file over to your `$HOME` directory
```sh
$ mkdir ~/.k8sclusterupgradetool/
$ cd k8sclusterupgradetool/
$ cp config.sample.yaml ~/.k8sclusterupgradetool/config.yaml
# make changes to the above file based on the versions of the components you want to check for the cluster
```

### Usage

Download the binary of the latest release from [Here](https://github.com/deliveryhero/k8s-cluster-upgrade-tool/releases)

On macOS, you might need to whitelist the binary to be able to run it using the following command:

```
sudo xattr -r -d com.apple.quarantine ~/path_to_binary/k8sclusterupgradetool
```

#### Running post upgrade checks

```
$ .k8sclusterupgradetool component version check -c=foo-cluster
2022/02/10 13:50:19 Please pass a valid clusterName

$ ./k8sclusterupgradetool component version check -c=valid-cluster-name
2022/03/25 13:44:15 Config file used: /Users/t.rahman/.k8sclusterupgradetool/config.yaml
2022/03/25 13:44:15 aws-node version read from config: aws-component-version
2022/03/25 13:44:15 coredns version read from config: coredns-component-version
2022/03/25 13:44:15 kube-proxy version read from config: kube-proxy-component-version
2022/03/25 13:44:15 cluster-autoscaler version read from config: cluster-autoscaler-component-version
Setting kubernetes context to valid-cluster-name
running post upgrade checks
Checking aws-node version
AWS Node Version on aws-component-version ✓
Checking kube-proxy version
kube-proxy needs to be updated, is currently on foo-version, desired version: kube-proxy-component-version
Checking core-dns version
core-dns needs to be updated, is currently on baz-version, desired version: coredns-component-version
Checking cluster-autoscaler version
cluster-autoscaler needs to be updated, is currently on far-version, desired version: cluster-autoscaler-component-version
```

#### Setting component versions for outdated components
```
$ ./k8sclusterupgradetool component version set -c=valid-cluster-name -k8s-comp=coredns -v=coredns-component-version
Setting kubernetes context to valid-cluster-name
2022/02/10 12:41:06 coredns has been set to coredns-component-version in cluster

$ ./k8sclusterupgradetool component version set -c=valid-cluster-name -k8s-comp=aws-node -v=aws-component-version
Setting kubernetes context to valid-cluster-name
2022/02/10 12:39:49 aws-node has been set to aws-component-version in cluster

$ ./k8sclusterupgradetool component version set -c=valid-cluster-name -k8s-comp=aws-node -v=aws-component-version123asd
2022/03/25 13:41:55 Config file used: /Users/t.rahman/.k8sclusterupgradetool/config.yaml
2022/03/25 13:41:55 aws-node version read from config: aws-component-version
2022/03/25 13:41:55 coredns version read from config: coredns-component-version
2022/03/25 13:41:55 kube-proxy version read from config: kube-proxy-component-version
2022/03/25 13:41:55 cluster-autoscaler version read from config: cluster-autoscaler-component-version
2022/03/25 13:41:55 aws-node component version passed doesn't match the version in config, please check the value in config file

$ ./k8sclusterupgradetool component version set -c=valid-cluster-name -k8s-comp=foo-deployment -v=vfoo-wrong-version
2022/03/25 13:42:52 Config file used: /Users/t.rahman/.k8sclusterupgradetool/config.yaml
2022/03/25 13:42:52 aws-node version read from config: aws-component-version
2022/03/25 13:42:52 coredns version read from config: coredns-component-version
2022/03/25 13:42:52 kube-proxy version read from config: kube-proxy-component-version
2022/03/25 13:42:52 cluster-autoscaler version read from config: cluster-autoscaler-component-version
2022/03/25 13:42:52 please pass a valid component name from this list [coredns, cluster-autoscaler, kube-proxy, aws-node]
```

#### Taint and drain nodes

**NOTE** as a side effect of this command, the tool also modifies size of the max instance size of the ASG to be set to current desired instance count to prevent the ASG being drained to scale up during the upgrade process.

##### With dry mode on (default set to true)

```
$ ./k8sclusterupgradetool asg taint-and-drain -c=valid-cluster-name -a=valid-asg-hash
2022/02/16 23:54:08 Setting kubernetes context to valid-cluster-name02
2022/02/16 23:54:09 Running cordon and drain command in dry mode
2022/02/16 23:54:09 Instances which are going to be tainted and drained from the ASG passed
2022/02/16 23:54:09 {"InstanceId":"i-foo","PrivateDNS":"ip-foo-ip.eu-west-1.compute.internal","AsgName":"valid-asg-hash"}
2022/02/16 23:54:09 {"InstanceId":"i-baz","PrivateDNS":"ip-baz.eu-west-1.compute.internal","AsgName":"valid-asg-hash"}
2022/02/16 23:54:09 {"InstanceId":"i-far","PrivateDNS":"ip-far.eu-west-1.compute.internal","AsgName":"valid-asg-hash"}
```

##### With dry mode on set to false

```
$ ./k8sclusterupgradetool asg taint-and-drain -c=valid-cluster-name -a=valid-cluster-name --dry-run=false
2022/02/16 23:54:29 Setting kubernetes context to valid-cluster-name
2022/02/16 23:54:30 Running cordon and drain command in non-dry mode
2022/02/16 23:54:31 Instances which are going to be tainted and drained from the ASG passed
2022/02/16 23:54:31 {"InstanceId":"i-foo","PrivateDNS":"ip-foo-ip.eu-west-1.compute.internal","AsgName":"valid-cluster-name"}
2022/02/16 23:54:31 {"InstanceId":"i-baz","PrivateDNS":"ip-baz.eu-west-1.compute.internal","AsgName":"valid-cluster-name"}
2022/02/16 23:54:31 {"InstanceId":"i-far","PrivateDNS":"ip-far.eu-west-1.compute.internal","AsgName":"valid-cluster-name"}
2022/02/16 23:54:31 Tainting node: ip-foo-ip.eu-west-1.compute.internal
2022/02/16 23:54:33 taint output:
 node/ip-foo-ip.eu-west-1.compute.internal tainted
2022/02/16 23:54:33 Tainting node: ip-baz.eu-west-1.compute.internal
2022/02/16 23:54:34 taint output:
 node/ip-baz.eu-west-1.compute.internal tainted
2022/02/16 23:54:34 Tainting node: ip-far.eu-west-1.compute.internal
2022/02/16 23:54:36 taint output:
 node/ip-far.eu-west-1.compute.internal tainted
2022/02/16 23:54:36 Draining node: ip-foo-ip.eu-west-1.compute.internal
2022/02/16 23:54:38 drain output:
 node/ip-foo-ip.eu-west-1.compute.internal cordoned
node/ip-foo-ip.eu-west-1.compute.internal drained
2022/02/16 23:54:38 Draining node: ip-baz.eu-west-1.compute.internal
2022/02/16 23:55:21 drain output:
 node/ip-baz.eu-west-1.compute.internal cordoned
evicting pod default/busybox-sleep
pod/busybox-sleep evicted
node/ip-baz.eu-west-1.compute.internal drained
2022/02/16 23:55:21 Draining node: ip-far.eu-west-1.compute.internal
2022/02/16 23:55:24 drain output:
 node/ip-far.eu-west-1.compute.internal cordoned
node/ip-far.eu-west-1.compute.internal drained
```

## Dev setup

- Install go 1.17

## Tests

- Running test suite
```
$ go test ./... -v
```

## Linting

```
$ docker run --rm -v $(pwd):/app -w /app golangci/golangci-lint:v1.45.2 golangci-lint run -v
```

## Adding a new release

Check [RELEASE.md](https://github.com/deliveryhero/k8s-cluster-upgrade-tool/tree/master/docs/RELEASE.md)
