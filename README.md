A plugin to Drone Helm push plugin to an OCI Docker registry.

# Usage

The following settings changes this plugin's behavior.

* `chart_destination` temporary directory where chart archive is written (default `.packaged_charts`)
* `chart_path` directory containing the helm chart (default `.`)
* `registry_namespace` registry namespace under which the chart will be published
* `registry_password` registry password/token
* `registry_url` registry where the packaged chart will be published (default `registry-1.docker.io`)
* `registry_username` registry username

Below is an example Harness CI pipeline step that uses this plugin.

```yaml
             - step:
                  type: Plugin
                  name: push helm chart
                  identifier: push_helm_chart
                  spec:
                    connectorRef: account.harnessImage
                    image: jimsheldon/drone-helm3-push
                    settings:
                      chart_path: charts/drone
                      registry_url: registry.hub.docker.com
                      registry_username: jimsheldon
                      registry_password: <+secrets.getValue("docker_pat")>
                      registry_namespace: jimsheldon
```
