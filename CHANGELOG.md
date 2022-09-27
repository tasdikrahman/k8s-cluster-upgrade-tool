## CHANGELOG

### v0.4.1

#### Changes

- improve error handling for `component version set` and `component version check` flow, remove repetitive function calls.

### v0.4.0

#### Adds

- client-go for the commands flow for `component version check` and `component version set`, which were earlier using
  the command line outputs parsed to further proceed, which was prone to errors as reported by some users.
- e2e tests for `component version set` and it's verification for components aws-node, coredns, cluster-autoscaler, kube-proxy.

#### Changes

- removes the functions `KubectlGetImageCommand`, `KubectlSetImageCommand` as they are not used anymore with their interactions
  replaced with client-go interactions

### v0.3.0

#### Adds

- commands `asg` parent command along with `taint-and-drain` being the child subcommand, which replaces,
  `taint-and-drain-asg`.
- commands `component`, `version` subcommand, further nested with `set` and `check` parent subcommands which
replace `setComponentVersion` and `postUpgradeCheck`.

#### Breaking change

- binary names changes to `k8sclusterupgradetool`
- the commands `setComponentVersion`, `taint-and-drain-asg`, `postUpgradeCheck` have been deprecated
for the new commands.
- the location of the config location changes to `$HOME/.k8sclusterupgradetool`.

#### Changes

- error messages fixed for displaying exact error when the config is not being read properly by the CLI
- remove positional arguments and start taking flags for the arguments passed to component subcommands to
reduce the overload on the user to remember which argument number comes where.

### v0.2.0

#### Adds

- e2e tests for the command postUpgradeCheck

#### Breaking change
- the config keys being used, please refer the same config and rename the key attributes as provided in `config.sample.yaml`
  - Name -> ClusterName
  - type -> ObjectType
  - name -> DeploymentName
- Adds new keys for the config to check, namely, `ContainerName`, `Namespace`

#### Adds
- ability to read container image, namespace to be read from config for a component being updated, for cases when deployment name
is not the same as the container name for the container for which update is being made.

### v0.1.0

OSS release of the internal tool which we have
