# User Provided Service Broker

User Provided Service Broker is an example
[Open Service Broker](https://www.openservicebrokerapi.org/)
for use demonstrating the Kubernetes
Service Catalog.

For more information,
[visit the Service Catalog project on github](https://github.com/kubernetes-incubator/service-catalog).

## Installing the Chart

To install the chart with the release name `opensds-broker`:

```bash
$ helm install charts/opensds-broker --name opensds-broker --namespace opensds-broker
```

## Uninstalling the Chart

To uninstall/delete the `opensds-broker` deployment:

```bash
$ helm delete opensds-broker
```

The command removes all the Kubernetes components associated with the chart and
deletes the release.

## Configuration

The following tables lists the configurable parameters of the User Provided
Service Broker

| Parameter | Description | Default |
|-----------|-------------|---------|
| `image` | Image to use | `leonwanghui/opensds-broker:latest` |
| `imagePullPolicy` | `imagePullPolicy` for the ups-broker | `Always` |

Specify each parameter using the `--set key=value[,key=value]` argument to
`helm install`.

Alternatively, a YAML file that specifies the values for the parameters can be
provided while installing the chart. For example:

```bash
$ helm install charts/opensds-broker --name opensds-broker --namespace opensds-broker \
  --values values.yaml
```
