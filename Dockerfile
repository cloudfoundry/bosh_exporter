FROM        quay.io/prometheus/busybox:latest
MAINTAINER  Ferran Rodenas <frodenas@gmail.com>

COPY bosh_exporter /bin/bosh_exporter

ENTRYPOINT ["/bin/bosh_exporter"]
EXPOSE     9190
