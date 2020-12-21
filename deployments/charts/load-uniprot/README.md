# load-uniprot

This [helm](https://github.com/kubernetes/helm) chart provides [kubernetes](http://kubernetes.io) manifests for running a [load-uniprot](https://hub.docker.com/r/dictybase/arangoadmin/) job.

# Managing the chart

## Install

```
helm install --name dev-release load-uniprot
```

For details, look [here](https://docs.helm.sh/using_helm/#helm-install-installing-a-package).

## Uninstall

```
helm uninstall dev-release
```

For details, look [here](https://docs.helm.sh/using_helm/#uninstall-a-release).

For upgrades and rollback, look [here](https://docs.helm.sh/using_helm/#helm-upgrade-and-helm-rollback-upgrading-a-release-and-recovering-on-failure).

## Configuration

The following tables lists the configurable parameters of the **chado-sqitch** chart and their default values.

| Parameter          | Description                                 | Default                 |
| ------------------ | ------------------------------------------- | ----------------------- |
| `image.repository` | load-uniprot image                          | `dictybase/arangoadmin` |
| `image.tag`        | image tag                                   | `0.0.1`                 |
| `image.pullPolicy` | Image pull policy                           | `IfNotPresent`          |
| `log.level`        | Name of user with ArangoDB admin privileges | `info`                  |

Specify each parameter using the `--set key=value[,key=value]` argument to `helm install`.

Alternatively, a YAML file that specifies the values for the parameters can be provided while installing the chart. For example,

```bash
$ helm install --name my-release -f values.yaml load-uniprot
```
