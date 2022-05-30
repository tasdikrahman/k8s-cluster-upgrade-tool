## RELEASE

# Before releasing a new version

Please complete all E2E tests as described in the docs.

# Release a new version

## Tools used

- [goreleaser](https://goreleaser.com/)
- git tags
- Drone pipelines

## Binaries released

- `k8s-cluster-upgrade-tool`

configured [here](../.goreleaser.yaml) under `builds`

### To release a new version

1. Create a new tag on Github.

    ```sh
    git tag -a v0.1.0 -m "First release"
    git push origin v0.1.0
    ```

   Tags should adhere to [semantic versioning](https://goreleaser.com/limitations/semver/) according to goreleaser documentation.

2. Drone will create a new release using goreleaser.
3. Released assets will be available [here](https://github.com/deliveryhero/k8s-cluster-upgrade-tool/releases).
