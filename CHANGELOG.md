## CHANGELOG

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
