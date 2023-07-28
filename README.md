# BOSH Prometheus Exporter [![Build Status](https://travis-ci.org/bosh-prometheus/bosh_exporter.png)](https://travis-ci.org/bosh-prometheus/bosh_exporter)

A [Prometheus][prometheus] exporter for [BOSH][bosh] metrics. Please refer to the [FAQ][faq] for general questions about
this exporter.

## Architecture overview

![](https://cdn.rawgit.com/bosh-prometheus/bosh_exporter/master/architecture/architecture.svg)

## Installation

### Binaries

Download the already existing [binaries][binaries] for your platform:

```bash
$ ./bosh_exporter <flags>
```

### From source

Using the standard `go install` (you must have [Go][golang] already installed in your local machine):

```bash
$ go install github.com/bosh-prometheus/bosh_exporter
$ bosh_exporter <flags>
```

### Docker

To run the bosh exporter as a Docker container, run:

```bash
$ docker run -p 9190:9190 boshprometheus/bosh-exporter <flags>
```

### Cloud Foundry

The exporter can be deployed to an already existing [Cloud Foundry][cloudfoundry] environment:

```bash
$ git clone https://github.com/bosh-prometheus/bosh_exporter.git
$ cd bosh_exporter
```

Modify the included [application manifest file][manifest] to include your BOSH properties. Then you can push the
exporter to your Cloud Foundry environment:

```bash
$ cf push
```

### BOSH

This exporter can be deployed using the [Prometheus BOSH Release][prometheus-boshrelease].

## Usage

### Flags

| Flag / Environment Variable                                          | Required | Default                   | Description                                                                                                                                                                                                                           |
|----------------------------------------------------------------------|----------|---------------------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `bosh.url`<br />`BOSH_EXPORTER_BOSH_URL`                             | Yes      |                           | BOSH URL                                                                                                                                                                                                                              |
| `bosh.username`<br />`BOSH_EXPORTER_BOSH_USERNAME`                   | *[1]*    |                           | BOSH Username                                                                                                                                                                                                                         |
| `bosh.password`<br />`BOSH_EXPORTER_BOSH_PASSWORD`                   | *[1]*    |                           | BOSH Password                                                                                                                                                                                                                         |
| `bosh.uaa.client-id`<br />`BOSH_EXPORTER_BOSH_UAA_CLIENT_ID`         | *[1]*    |                           | BOSH UAA Client ID                                                                                                                                                                                                                    |
| `bosh.uaa.client-secret`<br />`BOSH_EXPORTER_BOSH_UAA_CLIENT_SECRET` | *[1]*    |                           | BOSH UAA Client Secret                                                                                                                                                                                                                |
| `bosh.log-level`<br />`BOSH_EXPORTER_BOSH_LOG_LEVEL`                 | No       | `ERROR`                   | BOSH Log Level (`DEBUG`, `INFO`, `WARN`, `ERROR`, `NONE`)                                                                                                                                                                             |
| `bosh.ca-cert-file`<br />`BOSH_EXPORTER_BOSH_CA_CERT_FILE`           | Yes      |                           | BOSH CA Certificate file                                                                                                                                                                                                              |
| `filter.deployments`<br />`BOSH_EXPORTER_FILTER_DEPLOYMENTS`         | No       |                           | Comma separated deployments to filter                                                                                                                                                                                                 |
| `filter.azs`<br />`BOSH_EXPORTER_FILTER_AZS`                         | No       |                           | Comma separated AZs to filter                                                                                                                                                                                                         |
| `filter.collectors`<br />`BOSH_EXPORTER_FILTER_COLLECTORS`           | No       |                           | Comma separated collectors to filter. If not set, all collectors will be enabled  (`Deployments`, `Jobs`, `ServiceDiscovery`)                                                                                                         |
| `filter.cidrs`<br />`BOSH_EXPORTER_FILTER_CIDRS`                     | No       | `0.0.0.0/0`               | Comma separated CIDR to filter instance IPs                                                                                                                                                                                           |
| `metrics.namespace`<br />`BOSH_EXPORTER_METRICS_NAMESPACE`           | No       | `bosh`                    | Metrics Namespace                                                                                                                                                                                                                     |
| `metrics.environment`<br />`BOSH_EXPORTER_METRICS_ENVIRONMENT`       | Yes      |                           | Environment label to be attached to metrics                                                                                                                                                                                           |
| `sd.filename`<br />`BOSH_EXPORTER_SD_FILENAME`                       | No       | `bosh_target_groups.json` | Full path to the Service Discovery output file                                                                                                                                                                                        |
| `sd.processes_regexp`<br />`BOSH_EXPORTER_SD_PROCESSES_REGEXP`       | No       |                           | Regexp to filter Service Discovery processes names                                                                                                                                                                                    |
| `web.listen-address`<br />`BOSH_EXPORTER_WEB_LISTEN_ADDRESS`         | No       | `:9190`                   | Address to listen on for web interface and telemetry                                                                                                                                                                                  |
| `web.telemetry-path`<br />`BOSH_EXPORTER_WEB_TELEMETRY_PATH`         | No       | `/metrics`                | Path under which to expose Prometheus metrics                                                                                                                                                                                         |
| `web.auth.username`<br />`BOSH_EXPORTER_WEB_AUTH_USERNAME`           | No       |                           | Username for web interface basic auth                                                                                                                                                                                                 |
| `web.auth.password`<br />`BOSH_EXPORTER_WEB_AUTH_PASSWORD`           | No       |                           | Password for web interface basic auth                                                                                                                                                                                                 |
| `web.tls.cert_file`<br />`BOSH_EXPORTER_WEB_TLS_CERTFILE`            | No       |                           | Path to a file that contains the TLS certificate (PEM format). If the certificate is signed by a certificate authority, the file should be the concatenation of the server's certificate, any intermediates, and the CA's certificate |
| `web.tls.key_file`<br />`BOSH_EXPORTER_WEB_TLS_KEYFILE`              | No       |                           | Path to a file that contains the TLS private key (PEM format)                                                                                                                                                                         |

*[1]* When BOSH delegates user managament to [UAA][bosh_uaa], either `bosh.username` and `bosh.password`
or `bosh.uaa.client-id` and `bosh.uaa.client-secret` flags may be used; otherwise `bosh.username` and `bosh.password`
will be required. When using [UAA][bosh_uaa] and the `bosh.username` and `bosh.password` authentication method, tokens
are not refreshed, so after a period of time the exporter will be unable to communicate with the BOSH API, so use this
method only when testing the exporter. For production, it is recommended to use the `bosh.uaa.client-id`
and `bosh.uaa.client-secret` authentication method.

### Metrics

The exporter returns the following metrics:

| Metric                                               | Description                                                                                        | Labels                                  |
|------------------------------------------------------|----------------------------------------------------------------------------------------------------|-----------------------------------------|
| *metrics.namespace*\_scrapes\_total                  | Total number of times BOSH was scraped for metrics                                                 | `environment`, `bosh_name`, `bosh_uuid` |
| *metrics.namespace*\_scrape\_errors\_total           | Total number of times an error occured scraping BOSH                                               | `environment`, `bosh_name`, `bosh_uuid` |
| *metrics.namespace*\_last\_scrape\_error             | Whether the last scrape of metrics from BOSH resulted in an error (`1` for error, `0` for success) | `environment`, `bosh_name`, `bosh_uuid` |
| *metrics.namespace*\_last\_scrape\_timestamp         | Number of seconds since 1970 since last scrape from BOSH                                           | `environment`, `bosh_name`, `bosh_uuid` |
| *metrics.namespace*\_last\_scrape\_duration\_seconds | Duration of the last scrape from BOSH                                                              | `environment`, `bosh_name`, `bosh_uuid` |

The exporter returns the following `Deployments` metrics:

| Metric                                                            | Description                                                                     | Labels                                                                                                                               |
|-------------------------------------------------------------------|---------------------------------------------------------------------------------|--------------------------------------------------------------------------------------------------------------------------------------|
| *metrics.namespace*\_deployment\_release\_info                    | Labeled BOSH Deployment Release Info with a constant `1` value                  | `environment`, `bosh_name`, `bosh_uuid`, `bosh_deployment`, `bosh_release_name`, `bosh_release_version`                              |
| *metrics.namespace*\_deployment\_release\_job\_info               | Labeled BOSH Deployment Release Job Info with a constant `1` value              | `environment`, `bosh_name`, `bosh_uuid`, `bosh_deployment`, `bosh_release_name`, `bosh_release_version`, `bosh_release_job_name`     |
| *metrics.namespace*\_deployment\_release\_package\_info           | Labeled BOSH Deployment Release Package Info with a constant `1` value          | `environment`, `bosh_name`, `bosh_uuid`, `bosh_deployment`, `bosh_release_name`, `bosh_release_version`, `bosh_release_package_name` |
| *metrics.namespace*\_deployment\_stemcell\_info                   | Labeled BOSH Deployment Stemcell Info with a constant `1` value                 | `environment`, `bosh_name`, `bosh_uuid`, `bosh_deployment`, `bosh_stemcell_name`, `bosh_stemcell_version`, `bosh_stemcell_os_name`   |
| *metrics.namespace*\_deployment\_instances                        | Number of instances in the deployment                                           | `environment`, `bosh_name`, `bosh_uuid`, `bosh_deployment`, `bosh_vm_type`                                                           |
| *metrics.namespace*\_last\_deployments\_scrape\_timestamp         | Number of seconds since 1970 since last scrape of Deployments metrics from BOSH | `environment`, `bosh_name`, `bosh_uuid`                                                                                              |
| *metrics.namespace*\_last\_deployments\_scrape\_duration\_seconds | Duration of the last scrape of Deployments metrics from BOSH                    | `environment`, `bosh_name`, `bosh_uuid`                                                                                              |

The exporter returns the following `Jobs` metrics:

| Metric                                                     | Description                                                                                                                 | Labels                                                                                                                                                                                                                                   |
|------------------------------------------------------------|-----------------------------------------------------------------------------------------------------------------------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| *metrics.namespace*\_job\_healthy                          | BOSH Job Healthy (1 for healthy, 0 for unhealthy)                                                                           | `environment`, `bosh_name`, `bosh_uuid`, `bosh_deployment`, `bosh_job_name`, `bosh_job_id`, `bosh_job_index`, `bosh_job_az`, `bosh_job_ip`                                                                                               |
| *metrics.namespace*\_job\_load\_avg01                      | BOSH Job Load avg01                                                                                                         | `environment`, `bosh_name`, `bosh_uuid`, `bosh_deployment`, `bosh_job_name`, `bosh_job_id`, `bosh_job_index`, `bosh_job_az`, `bosh_job_ip`                                                                                               |
| *metrics.namespace*\_job\_load\_avg05                      | BOSH Job Load avg05                                                                                                         | `environment`, `bosh_name`, `bosh_uuid`, `bosh_deployment`, `bosh_job_name`, `bosh_job_id`, `bosh_job_index`, `bosh_job_az`, `bosh_job_ip`                                                                                               |
| *metrics.namespace*\_job\_load\_avg15                      | BOSH Job Load avg15                                                                                                         | `environment`, `bosh_name`, `bosh_uuid`, `bosh_deployment`, `bosh_job_name`, `bosh_job_id`, `bosh_job_index`, `bosh_job_az`, `bosh_job_ip`                                                                                               |
| *metrics.namespace*\_job\_cpu\_sys                         | BOSH Job CPU System                                                                                                         | `environment`, `bosh_name`, `bosh_uuid`, `bosh_deployment`, `bosh_job_name`, `bosh_job_id`, `bosh_job_index`, `bosh_job_az`, `bosh_job_ip`                                                                                               |
| *metrics.namespace*\_job\_cpu\_user                        | BOSH Job CPU User                                                                                                           | `environment`, `bosh_name`, `bosh_uuid`, `bosh_deployment`, `bosh_job_name`, `bosh_job_id`, `bosh_job_index`, `bosh_job_az`, `bosh_job_ip`                                                                                               |
| *metrics.namespace*\_job\_cpu\_wait                        | BOSH Job CPU Wait                                                                                                           | `environment`, `bosh_name`, `bosh_uuid`, `bosh_deployment`, `bosh_job_name`, `bosh_job_id`, `bosh_job_index`, `bosh_job_az`, `bosh_job_ip`                                                                                               |
| *metrics.namespace*\_job\_mem\_kb                          | BOSH Job Memory KB                                                                                                          | `environment`, `bosh_name`, `bosh_uuid`, `bosh_deployment`, `bosh_job_name`, `bosh_job_id`, `bosh_job_index`, `bosh_job_az`, `bosh_job_ip`                                                                                               |
| *metrics.namespace*\_job\_mem\_percent                     | BOSH Job Memory Percent                                                                                                     | `environment`, `bosh_name`, `bosh_uuid`, `bosh_deployment`, `bosh_job_name`, `bosh_job_id`, `bosh_job_index`, `bosh_job_az`, `bosh_job_ip`                                                                                               |
| *metrics.namespace*\_job\_swap\_kb                         | BOSH Job Swap KB                                                                                                            | `environment`, `bosh_name`, `bosh_uuid`, `bosh_deployment`, `bosh_job_name`, `bosh_job_id`, `bosh_job_index`, `bosh_job_az`, `bosh_job_ip`                                                                                               |
| *metrics.namespace*\_job\_swap\_percent                    | BOSH Job Swap Percent                                                                                                       | `environment`, `bosh_name`, `bosh_uuid`, `bosh_deployment`, `bosh_job_name`, `bosh_job_id`, `bosh_job_index`, `bosh_job_az`, `bosh_job_ip`                                                                                               |
| *metrics.namespace*\_job\_system\_disk\_inode\_percent     | BOSH Job System Disk Inode Percent                                                                                          | `environment`, `bosh_name`, `bosh_uuid`, `bosh_deployment`, `bosh_job_name`, `bosh_job_id`, `bosh_job_index`, `bosh_job_az`, `bosh_job_ip`                                                                                               |
| *metrics.namespace*\_job\_system\_disk\_percent            | BOSH Job System Disk Percent                                                                                                | `environment`, `bosh_name`, `bosh_uuid`, `bosh_deployment`, `bosh_job_name`, `bosh_job_id`, `bosh_job_index`, `bosh_job_az`, `bosh_job_ip`                                                                                               |
| *metrics.namespace*\_job\_ephemeral\_disk\_inode\_percent  | BOSH Job Ephemeral Disk Inode Percent                                                                                       | `environment`, `bosh_name`, `bosh_uuid`, `bosh_deployment`, `bosh_job_name`, `bosh_job_id`, `bosh_job_index`, `bosh_job_az`, `bosh_job_ip`                                                                                               |
| *metrics.namespace*\_job\_ephemeral\_disk\_percent         | BOSH Job Ephemeral Disk Percent                                                                                             | `environment`, `bosh_name`, `bosh_uuid`, `bosh_deployment`, `bosh_job_name`, `bosh_job_id`, `bosh_job_index`, `bosh_job_az`, `bosh_job_ip`                                                                                               |
| *metrics.namespace*\_job\_persistent\_disk\_inode\_percent | BOSH Job Persistent Disk Inode Percent                                                                                      | `environment`, `bosh_name`, `bosh_uuid`, `bosh_deployment`, `bosh_job_name`, `bosh_job_id`, `bosh_job_index`, `bosh_job_az`, `bosh_job_ip`                                                                                               |
| *metrics.namespace*\_job\_persistent\_disk\_percent        | BOSH Job Persistent Disk Percent                                                                                            | `environment`, `bosh_name`, `bosh_uuid`, `bosh_deployment`, `bosh_job_name`, `bosh_job_id`, `bosh_job_index`, `bosh_job_az`, `bosh_job_ip`                                                                                               |
| *metrics.namespace*\_job\_process\_info                    | BOSH Job Process Info with a constant '1' value. Release can be found only if process name is the same as release job name. | `environment`, `bosh_name`, `bosh_uuid`, `bosh_deployment`, `bosh_job_name`, `bosh_job_id`, `bosh_job_index`, `bosh_job_az`, `bosh_job_ip`, `bosh_job_process_name`, `bosh_job_process_release_name`, `bosh_job_process_release_version` |
| *metrics.namespace*\_job\_process\_healthy                 | BOSH Job Process Healthy (1 for healthy, 0 for unhealthy)                                                                   | `environment`, `bosh_name`, `bosh_uuid`, `bosh_deployment`, `bosh_job_name`, `bosh_job_id`, `bosh_job_index`, `bosh_job_az`, `bosh_job_ip`, `bosh_job_process_name`                                                                      |
| *metrics.namespace*\_job\_process\_uptime\_seconds         | BOSH Job Process Uptime in seconds                                                                                          | `environment`, `bosh_name`, `bosh_uuid`, `bosh_deployment`, `bosh_job_name`, `bosh_job_id`, `bosh_job_index`, `bosh_job_az`, `bosh_job_ip`, `bosh_job_process_name`                                                                      |
| *metrics.namespace*\_job\_process\_cpu\_total              | BOSH Job Process CPU Total                                                                                                  | `environment`, `bosh_name`, `bosh_uuid`, `bosh_deployment`, `bosh_job_name`, `bosh_job_id`, `bosh_job_index`, `bosh_job_az`, `bosh_job_ip`, `bosh_job_process_name`                                                                      |
| *metrics.namespace*\_job\_process\_mem\_kb                 | BOSH Job Process Memory KB                                                                                                  | `environment`, `bosh_name`, `bosh_uuid`, `bosh_deployment`, `bosh_job_name`, `bosh_job_id`, `bosh_job_index`, `bosh_job_az`, `bosh_job_ip`, `bosh_job_process_name`                                                                      |
| *metrics.namespace*\_job\_process\_mem\_percent            | BOSH Job Process Memory Percent                                                                                             | `environment`, `bosh_name`, `bosh_uuid`, `bosh_deployment`, `bosh_job_name`, `bosh_job_id`, `bosh_job_index`, `bosh_job_az`, `bosh_job_ip`, `bosh_job_process_name`                                                                      |
| *metrics.namespace*\_last\_jobs\_scrape\_timestamp         | Number of seconds since 1970 since last scrape of Job metrics from BOSH                                                     | `environment`, `bosh_name`, `bosh_uuid`                                                                                                                                                                                                  |
| *metrics.namespace*\_last\_jobs\_scrape\_duration\_seconds | Duration of the last scrape of Job metrics from BOSH                                                                        | `environment`, `bosh_name`, `bosh_uuid`                                                                                                                                                                                                  |

The exporter returns the following `ServiceDiscovery` metrics:

| Metric                                                                   | Description                                                                   | Labels                                  |
|--------------------------------------------------------------------------|-------------------------------------------------------------------------------|-----------------------------------------|
| *metrics.namespace*\_last\_service\_discovery\_scrape\_timestamp         | Number of seconds since 1970 since last scrape of Service Discovery from BOSH | `environment`, `bosh_name`, `bosh_uuid` |
| *metrics.namespace*\_last\_service\_discovery\_scrape\_duration\_seconds | Duration of the last scrape of Service Discovery from BOSH                    | `environment`, `bosh_name`, `bosh_uuid` |

### Service Discovery

If the `ServiceDiscovery` collector is enabled, the exporter will write a `json` file at the `sd.filename` location
containing a list of static configs that can be used with the Prometheus [file-based service discovery][file_sd_config]
mechanism:

```json
[
  {
    "targets": [
      "10.244.0.12"
    ],
    "labels": {
      "__meta_bosh_job_process_name": "bosh_exporter"
    }
  },
  {
    "targets": [
      "10.244.0.11",
      "10.244.0.12",
      "10.244.0.13",
      "10.244.0.14"
    ],
    "labels": {
      "__meta_bosh_deployment": "deployment1",
      "__meta_bosh_deployment_releases": "exporters_release:1.0,other_release:0.2",
      "__meta_bosh_job_process_name": "node_exporter",
      "__meta_bosh_job_process_release":"exporters_release:1.0"
    }
  }
]
```

[!NOTE]
`__meta_bosh_job_process_release` has the same value as the labels: `bosh_job_process_release_name`:`bosh_job_process_release_version`.
The `process release` can be found only if the `process name` (label `bosh_job_process_name`) is the same as `BOSH release job name`.
`BOSH release job name` is not the same as the label `bosh_job_name` which is the `instance group name` (BOSH `deployment manifest`).


The list of targets can be filtered using the `sd.processes_regexp` flag.

[Prometheus file-based service discovery](https://prometheus.io/docs/prometheus/latest/configuration/configuration/#file_sd_config) example:
```yaml
- job_name: node_exporter
  metrics_path: /metrics
  scheme: http
  file_sd_configs:
    - files:
        - /var/vcap/store/bosh_exporter/bosh_target_groups.json
  relabel_configs:
    - source_labels: [ __meta_bosh_job_process_name ]
      regex: 'node_exporter'
      action: keep
    - source_labels: [ __meta_bosh_deployment_releases ]
      regex: '.*exporters_release:1\..*'
      action: keep

```

### Filtering IPs

Available instance IPs can be filtered using the `filter.cidrs` flag.

The first IP that matches a CIDR is used as target. CIDRs are tested in the order specified by the comma-seperated list.
The instance is dropped if no IP is included in any of the CIDRs.

## Contributing

Refer to the [contributing guidelines][contributing].

## License

Apache License 2.0, see [LICENSE][license].

[binaries]: https://github.com/bosh-prometheus/bosh_exporter/releases

[bosh]: https://bosh.io

[bosh_uaa]: https://bosh.io/docs/director-users-uaa/

[cloudfoundry]: https://www.cloudfoundry.org/

[contributing]: https://github.com/bosh-prometheus/bosh_exporter/blob/master/CONTRIBUTING.md

[faq]: https://github.com/bosh-prometheus/bosh_exporter/blob/master/FAQ.md

[file_sd_config]: https://prometheus.io/docs/prometheus/latest/configuration/configuration/#file_sd_config

[golang]: https://go.dev/

[license]: https://github.com/bosh-prometheus/bosh_exporter/blob/master/LICENSE

[manifest]: https://github.com/bosh-prometheus/bosh_exporter/blob/master/manifest.yml

[prometheus]: https://prometheus.io/

[prometheus-boshrelease]: https://github.com/bosh-prometheus/prometheus-boshrelease
