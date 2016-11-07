# FAQ

### What metrics does this exporter report?

The BOSH Prometheus Exporter gets the metrics from the [BOSH Director][bosh_director], who gathers them from each VM [BOSH Agent][bosh_agent]. The metrics that are being [reported][bosh_exporter_metrics] are pretty basic, but include:

* Deployment metrics:
  * Releases in use
  * Stemcells in use
* Job metrics:
  * Health status
  * CPU
  * Load
  * Memory
  * Swap
  * System, Ephemeral and Persistent disk
* Process metrics:
  * Health status
  * Uptime
  * CPU
  * Memory

### How can I get more detailed metrics from each VM?

If you want to get more detailed VM system metrics, like disk I/O, network traffic, ..., it is recommended to deploy the Prometheus [Node exporter][node_exporter] on each VM.

### What are the caveats when using this exporter?

In order to get the metrics, the exporter calls the [BOSH Director][bosh_director] [Instance details][instance_details] endpoint. This request result in potentially long running operations against each [BOSH Agent][bosh_agent], so such requests start a [Director Task][director_task]. Therefore, each exporter scrape will generate a new [Director Task][director_task] per deployment. This will NOT hurt your BOSH performance, but has the nasty effect that generates thousand of tasks per scrape (i.e. scrapping each minute will generate 1440 tasks per deployment per day).

It is, therefore, recommended to increase the `scrape interval` and the `scrape timeout` for this exporter:

```yaml
scrape_configs:
  - job_name: bosh_exporter
    scrape_interval: 2m
    scrape_timeout: 1m
```

A longer `scrape interval` means less *real time* metrics, but for most use cases, this will be enough, specially when combined with the [Node exporter][node_exporter].

### How can I get BOSH metrics without the above caveats?

An alternative approach to gather BOSH metrics without using this exporter is to use the [Graphite exporter][graphite_exporter] and configure a [metric mapping][graphite_mapping]:

```
*.*.*.*.system_healthy
name="bosh_exporter_job_healthy"
bosh_deployment="$1"
bosh_job_name="$2"
bosh_job_id="$3"
...
```

Then you will need to enable the [Graphite plugin][bosh_graphite] at your [BOSH Health Monitor][bosh_health_monitor] configuration pointing to the [Graphite exporter][graphite_exporter] IP address.

On the other hand, the downside of this approach is that you will NOT get the same level of metrics that this exporter reports and you cannot use the service discovery approach.

### How can I enable only a particular collector?

The *bosh.collectors* command flag allows you to filter what collectors will be enabled. Possible values are `Deployments` ,`Jobs`, `ServiceDiscovery` (or a combination of them).

### How can I filter by a particular BOSH deployment?

The *bosh.deployments* command flag allows you to filter what BOSH deployments will be reported.

### How can I get the BOSH CA certificate?

Comunication between the exporter and the [BOSH Director][bosh_director] uses HTTPS. Actually, there is no way to disable the SSL certificate validation, so therefore, the certificates must be created setting a [Subject Alternative Name][san] (SAN) with the IP address of the [BOSH Director][bosh_director], otherwise, you will get the following error message:

```
x509: cannot validate certificate for X.X.X.X because it doesn't contain any IP SANs
```

In order to generate the proper certificates, you can reuse the [generate.sh][generate_certificates] script located at the [BOSH Lite][bosh_lite] repository updating the IP address. Then deploy (or redeploy) your [BOSH Director][bosh_director] adding the following properties:

```yaml
properties:
  director:
    ssl:
      cert: |
        <content of the director.crt file>
      key: |
        <content of the director.key file>
```

Later, when starting the `bosh_exporter` you must specify the *bosh.ca-cert-file* command line flag pointing to the location of the `ca.crt` file.

For testing purposes, this repository includes the [CA Cert][bosh_lite_ca_cert] to be used only when testing the exporter against a [BOSH Lite][bosh_lite].

### How can I use the Service Discovery?

If you don't want to configure manually all exporters IP addresses at your prometheus configuration file, you can use the Prometheus [file-based service discovery][file_sd_config] mechanism. Just point the `file_sd_configs` configuration to the output file (*sd.filename* command flag) of this exporter and use the Prometheus [relabel configuration][relabel_config] to get the IP address:

```yaml
scrape_configs:
  - job_name: node_exporter
    file_sd_configs:
      - files:
        - /var/vcap/store/bosh_exporter/bosh_target_groups.json
    relabel_configs:
      - source_labels: [__meta_bosh_job_process_name]
        regex: node_exporter
        action: keep
      - source_labels: [__address__]
        regex: "(.*)"
        target_label: __address__
        replacement: "${1}:9100"
```

### How can I filter the Service Discovery output file by a particular exporter?

Yes, the *sd.processes_regexp* command flag allows you to filter what BOSH Job processes will be reported.

### Why is the BOSH Service Discovery a collector?

There are mainly two reasons:

* Prometheus Service Discovery is not pluggable, that means that you either incorporate the BOSH Service Discovery as part of the official [Prometheus core code][prometheus_github] or you create a separate executable that produces an output file that can be used by the Prometheus [file-based service discovery][file_sd_config] mechanism. We decided to use the [file-based service discovery][file_sd_config] mechanism because it was easier for us to test this approach.
* We want to minimize the number of calls to the [BOSH Director][bosh_director] (see the above [caveats](#how-can-i-get-bosh-metrics-without-the-above-caveats)). Having a different executable means that in order to get the BOSH Job IPs and processes we will need to generate a new [Director Task][director_task]. Using a new collector within this exporter allows us to reuse the same deployment calls.

### What is the recommended deployment strategy?

Prometheus advises to collocate exporters near the metrics source, in this case, that means colocating this exporter within your [BOSH Director][bosh_director] VM. We encourage you to follow this approach whenever is possible.

But the downside of the above advice is when using the Service Discovery mechanism. In this case, the exporter must be located at the Prometheus VM in order to access the service discovery output file.

### I have a question but I don't see it answered at this FAQ

We will be glad to address any questions not answered here. Please, just open a [new issue][issues].

[bosh_agent]: https://bosh.io/docs/bosh-components.html#agent
[bosh_director]: http://bosh.io/docs/bosh-components.html#director
[bosh_exporter_metrics]: https://github.com/cloudfoundry-community/bosh_exporter#metrics
[bosh_graphite]: http://bosh.io/docs/hm-config.html#graphite
[bosh_health_monitor]: http://bosh.io/docs/bosh-components.html#health-monitor
[bosh_lite]: https://github.com/cloudfoundry/bosh-lite
[bosh_lite_ca_cert]: https://github.com/cloudfoundry-community/bosh_exporter/blob/master/bosh-lite-ca.crt
[director_task]: http://bosh.io/docs/director-tasks.html
[file_sd_config]: https://prometheus.io/docs/operating/configuration/#&lt;file_sd_config&gt;
[graphite_exporter]: https://github.com/prometheus/graphite_exporter
[graphite_mapping]: https://github.com/prometheus/graphite_exporter#metric-mapping-and-configuration
[instance_details]: https://bosh.io/docs/director-api-v1.html#list-instances-detailed
[issues]: https://github.com/cloudfoundry-community/bosh_exporter/issues
[node_exporter]: https://github.com/prometheus/node_exporter
[prometheus_github]: https://github.com/prometheus/prometheus
[relabel_config]: https://prometheus.io/docs/operating/configuration/#<relabel_config>
[san]: https://en.wikipedia.org/wiki/Subject_Alternative_Name
[generate_certificates]: https://github.com/cloudfoundry/bosh-lite/blob/master/ca/generate.sh
