# BOSH Exporter [![Build Status](https://travis-ci.org/cloudfoundry-community/bosh_exporter.png)](https://travis-ci.org/cloudfoundry-community/bosh_exporter)

A [Prometheus][prometheus] exporter for [BOSH][bosh] metrics. It exports the following metrics for every `deployment` and `process`:

| Metric | Description | Labels |
| ------ | ----------- | ------ |
| *namespace*_bosh_job_healthy | BOSH Job Healthy | bosh_deployment, bosh_job, bosh_index, bosh_az
| *namespace*_bosh_job_load_avg01 | BOSH Job Load avg01 | bosh_deployment, bosh_job, bosh_index, bosh_az
| *namespace*_bosh_job_load_avg05 | BOSH Job Load avg05 | bosh_deployment, bosh_job, bosh_index, bosh_az
| *namespace*_bosh_job_load_avg15 | BOSH Job Load avg15 | bosh_deployment, bosh_job, bosh_index, bosh_az
| *namespace*_bosh_job_cpu_system | BOSH Job CPU System | bosh_deployment, bosh_job, bosh_index, bosh_az
| *namespace*_bosh_job_cpu_user | BOSH Job CPU User | bosh_deployment, bosh_job, bosh_index, bosh_az
| *namespace*_bosh_job_cpu_wait | BOSH Job CPU Wait | bosh_deployment, bosh_job, bosh_index, bosh_az
| *namespace*_bosh_job_mem_kb | BOSH Job Memory KB | bosh_deployment, bosh_job, bosh_index, bosh_az
| *namespace*_bosh_job_mem_percent | BOSH Job Memory Percent | bosh_deployment, bosh_job, bosh_index, bosh_az
| *namespace*_bosh_job_swap_kb | BOSH Job Swap KB | bosh_deployment, bosh_job, bosh_index, bosh_az
| *namespace*_bosh_job_swap_percent | BOSH Job Swap Percent | bosh_deployment, bosh_job, bosh_index, bosh_az
| *namespace*_bosh_job_system_disk_inode_percent | BOSH Job System Disk Inode Percent | bosh_deployment, bosh_job, bosh_index, bosh_az
| *namespace*_bosh_job_system_disk_percent | BOSH Job System Disk Percent | bosh_deployment, bosh_job, bosh_index, bosh_az
| *namespace*_bosh_job_ephemeral_disk_inode_percent | BOSH Job Ephemeral Disk Inode Percent | bosh_deployment, bosh_job, bosh_index, bosh_az
| *namespace*_bosh_job_ephemeral_disk_percent | BOSH Job Ephemeral Disk Percent | bosh_deployment, bosh_job, bosh_index, bosh_az
| *namespace*_bosh_job_persistent_disk_inode_percent | BOSH Job Persistent Disk Inode Percent | bosh_deployment, bosh_job, bosh_index, bosh_az
| *namespace*_bosh_job_persistent_disk_percent | BOSH Job Persistent Disk Percent | bosh_deployment, bosh_job, bosh_index, bosh_az
| *namespace*_bosh_job_process_healthy | BOSH Job Process Healthy | bosh_deployment, bosh_job, bosh_index, bosh_az, bosh_process_name
| *namespace*_bosh_job_process_uptime_seconds | BOSH Job Process Uptime in seconds | bosh_deployment, bosh_job, bosh_index, bosh_az, bosh_process_name
| *namespace*_bosh_job_process_cpu_total | BOSH Job Process CPU Total | bosh_deployment, bosh_job, bosh_index, bosh_az, bosh_process_name
| *namespace*_bosh_job_process_mem_kb | BOSH Job Process Memory KB | bosh_deployment, bosh_job, bosh_index, bosh_az, bosh_process_name
| *namespace*_bosh_job_process_mem_percent | BOSH Job Process Memory Percent | bosh_deployment, bosh_job, bosh_index, bosh_az, bosh_process_name

## Building and running

```bash
make
./bosh_exporter <flags>
```

### Flags

| Flag / Environment Variable | Required | Default | Description
| --------------------------- | -------- | ------- | -----------
| bosh.url<br />BOSH_EXPORTER_BOSH_URL | Yes | | BOSH URL
| bosh.username<br />BOSH_EXPORTER_BOSH_USERNAME | No | | BOSH Username
| bosh.password<br />BOSH_EXPORTER_BOSH_PASSWORD | No | | BOSH Password
| bosh.log-level<br />BOSH_EXPORTER_BOSH_LOG_LEVEL | No | ERROR | BOSH Log Level ("DEBUG", "INFO", "WARN", "ERROR", "NONE")
| bosh.ca-cert-file<br />BOSH_EXPORTER_BOSH_CA_CERT_FILE | No | | BOSH CA Certificate file
| uaa.url<br />BOSH_EXPORTER_UAA_URL | No | | BOSH UAA URL
| uaa.client-id<br />BOSH_EXPORTER_UAA_CLIENT_ID | No | | BOSH UAA Client ID
| uaa.client-secret<br />BOSH_EXPORTER_UAA_CLIENT_SECRET | No | | BOSH UAA Client Secret
| metrics.namespace<br />BOSH_EXPORTER_METRICS_NAMESPACE | No | bosh_exporter | Metrics Namespace
| web.listen-address<br />BOSH_EXPORTER_WEB_LISTEN_ADDRESS | No | :9190 | Address to listen on for web interface and telemetry
| web.telemetry-path<br />BOSH_EXPORTER_WEB_TELEMETRY_PATH | No | /metrics | Path under which to expose Prometheus metrics

## Running tests

```bash
make test
```

## Using Docker

You can deploy this exporter using the [frodenas/bosh-exporter][hub] Docker image. For example:

```bash
docker pull frodenas/bosh-exporter

docker run -d -p 9190:9190 frodenas/bosh-exporter <flags>
```

[bosh]: https://bosh.io
[hub]: https://hub.docker.com/r/frodenas/bosh-exporter/
[prometheus]: https://prometheus.io/
