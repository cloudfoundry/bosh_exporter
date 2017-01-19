package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/cloudfoundry/bosh-cli/director"
	"github.com/cloudfoundry/bosh-cli/uaa"
	"github.com/cloudfoundry/bosh-utils/logger"
	"github.com/cloudfoundry/bosh-utils/system"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"github.com/prometheus/common/version"

	"github.com/cloudfoundry-community/bosh_exporter/collectors"
	"github.com/cloudfoundry-community/bosh_exporter/deployments"
	"github.com/cloudfoundry-community/bosh_exporter/filters"
)

var (
	boshURL = flag.String(
		"bosh.url", "",
		"BOSH URL ($BOSH_EXPORTER_BOSH_URL).",
	)

	boshUsername = flag.String(
		"bosh.username", "",
		"BOSH Username ($BOSH_EXPORTER_BOSH_USERNAME).",
	)

	boshPassword = flag.String(
		"bosh.password", "",
		"BOSH Password ($BOSH_EXPORTER_BOSH_PASSWORD).",
	)

	boshUAAURL = flag.String(
		"bosh.uaa.url", "",
		"BOSH UAA Url ($BOSH_EXPORTER_BOSH_UAA_URL).",
	)

	boshUAAClientID = flag.String(
		"bosh.uaa.client-id", "",
		"BOSH UAA Client ID ($BOSH_EXPORTER_BOSH_UAA_CLIENT_ID).",
	)

	boshUAAClientSecret = flag.String(
		"bosh.uaa.client-secret", "",
		"BOSH UAA Client Secret ($BOSH_EXPORTER_BOSH_UAA_CLIENT_SECRET).",
	)

	boshLogLevel = flag.String(
		"bosh.log-level", "ERROR",
		"BOSH Log Level ($BOSH_EXPORTER_BOSH_LOG_LEVEL).",
	)

	boshCACertFile = flag.String(
		"bosh.ca-cert-file", "",
		"BOSH CA Certificate file ($BOSH_EXPORTER_BOSH_CA_CERT_FILE).",
	)

	filterDeployments = flag.String(
		"filter.deployments", "",
		"Comma separated deployments to filter ($BOSH_EXPORTER_FILTER_DEPLOYMENTS).",
	)

	filterAZs = flag.String(
		"filter.azs", "",
		"Comma separated AZs to filter ($BOSH_EXPORTER_FILTER_AZS).",
	)

	filterCollectors = flag.String(
		"filter.collectors", "",
		"Comma separated collectors to filter (Deployments,Jobs,ServiceDiscovery) ($BOSH_EXPORTER_FILTER_COLLECTORS).",
	)

	metricsNamespace = flag.String(
		"metrics.namespace", "bosh",
		"Metrics Namespace ($BOSH_EXPORTER_METRICS_NAMESPACE).",
	)

	sdFilename = flag.String(
		"sd.filename", "bosh_target_groups.json",
		"Full path to the Service Discovery output file ($BOSH_EXPORTER_SD_FILENAME).",
	)

	sdProcessesRegexp = flag.String(
		"sd.processes_regexp", "",
		"Regexp to filter Service Discovery processes names ($BOSH_EXPORTER_SD_PROCESSES_REGEXP).",
	)

	showVersion = flag.Bool(
		"version", false,
		"Print version information.",
	)

	listenAddress = flag.String(
		"web.listen-address", ":9190",
		"Address to listen on for web interface and telemetry ($BOSH_EXPORTER_WEB_LISTEN_ADDRESS).",
	)

	metricsPath = flag.String(
		"web.telemetry-path", "/metrics",
		"Path under which to expose Prometheus metrics ($BOSH_EXPORTER_WEB_TELEMETRY_PATH).",
	)
)

func init() {
	prometheus.MustRegister(version.NewCollector(*metricsNamespace))
}

func overrideFlagsWithEnvVars() {
	overrideWithEnvVar("BOSH_EXPORTER_BOSH_URL", boshURL)
	overrideWithEnvVar("BOSH_EXPORTER_BOSH_USERNAME", boshUsername)
	overrideWithEnvVar("BOSH_EXPORTER_BOSH_PASSWORD", boshPassword)
	overrideWithEnvVar("BOSH_EXPORTER_BOSH_UAA_URL", boshUAAURL)
	overrideWithEnvVar("BOSH_EXPORTER_BOSH_UAA_CLIENT_ID", boshUAAClientID)
	overrideWithEnvVar("BOSH_EXPORTER_BOSH_UAA_CLIENT_SECRET", boshUAAClientSecret)
	overrideWithEnvVar("BOSH_EXPORTER_BOSH_LOG_LEVEL", boshLogLevel)
	overrideWithEnvVar("BOSH_EXPORTER_BOSH_CA_CERT_FILE", boshCACertFile)
	overrideWithEnvVar("BOSH_EXPORTER_FILTER_DEPLOYMENTS", filterDeployments)
	overrideWithEnvVar("BOSH_EXPORTER_FILTER_AZS", filterAZs)
	overrideWithEnvVar("BOSH_EXPORTER_FILTER_COLLECTORS", filterCollectors)
	overrideWithEnvVar("BOSH_EXPORTER_METRICS_NAMESPACE", metricsNamespace)
	overrideWithEnvVar("BOSH_EXPORTER_SD_FILENAME", sdFilename)
	overrideWithEnvVar("BOSH_EXPORTER_SD_PROCESSES_REGEXP", sdProcessesRegexp)
	overrideWithEnvVar("BOSH_EXPORTER_WEB_LISTEN_ADDRESS", listenAddress)
	overrideWithEnvVar("BOSH_EXPORTER_WEB_TELEMETRY_PATH", metricsPath)
}

func overrideWithEnvVar(name string, value *string) {
	envValue := os.Getenv(name)
	if envValue != "" {
		*value = envValue
	}
}

func readCACert(CACertFile string, logger logger.Logger) (string, error) {
	if CACertFile != "" {
		fs := system.NewOsFileSystem(logger)

		CACertFileFullPath, err := fs.ExpandPath(CACertFile)
		if err != nil {
			return "", err
		}

		CACert, err := fs.ReadFileString(CACertFileFullPath)
		if err != nil {
			return "", err
		}

		return CACert, nil
	}

	return "", nil
}

func buildBOSHClient() (director.Director, error) {
	logLevel, err := logger.Levelify(*boshLogLevel)
	if err != nil {
		return nil, err
	}

	logger := logger.NewLogger(logLevel)

	directorConfig, err := director.NewConfigFromURL(*boshURL)
	if err != nil {
		return nil, err
	}

	boshCACert, err := readCACert(*boshCACertFile, logger)
	if err != nil {
		return nil, err
	}

	directorConfig.CACert = boshCACert
	directorConfig.Client = *boshUsername
	directorConfig.ClientSecret = *boshPassword

	if *boshUAAURL != "" {
		uaaConfig, err := uaa.NewConfigFromURL(*boshUAAURL)
		if err != nil {
			return nil, err
		}

		uaaConfig.CACert = boshCACert
		uaaConfig.Client = *boshUAAClientID
		uaaConfig.ClientSecret = *boshUAAClientSecret

		uaaFactory := uaa.NewFactory(logger)
		uaaClient, err := uaaFactory.New(uaaConfig)
		if err != nil {
			return nil, err
		}

		directorConfig.TokenFunc = uaa.NewClientTokenSession(uaaClient).TokenFunc
	}

	boshFactory := director.NewFactory(logger)
	boshClient, err := boshFactory.New(directorConfig, director.NewNoopTaskReporter(), director.NewNoopFileReporter())
	if err != nil {
		return nil, err
	}

	return boshClient, nil
}

func main() {
	flag.Parse()
	overrideFlagsWithEnvVars()

	if *showVersion {
		fmt.Fprintln(os.Stdout, version.Print("bosh_exporter"))
		os.Exit(0)
	}

	log.Infoln("Starting bosh_exporter", version.Info())
	log.Infoln("Build context", version.BuildContext())

	boshClient, err := buildBOSHClient()
	if err != nil {
		log.Errorf("Error creating BOSH Client: %s", err.Error())
		os.Exit(1)
	}

	boshInfo, err := boshClient.Info()
	if err != nil {
		log.Errorf("Error reading BOSH Info: %s", err.Error())
		os.Exit(1)
	}
	log.Infof("Using BOSH Director `%s` (%s)", boshInfo.Name, boshInfo.UUID)

	var deploymentsFilters []string
	if *filterDeployments != "" {
		deploymentsFilters = strings.Split(*filterDeployments, ",")
	}
	deploymentsFilter := filters.NewDeploymentsFilter(deploymentsFilters, boshClient)
	deploymentsFetcher := deployments.NewFetcher(*deploymentsFilter)

	var azsFilters []string
	if *filterAZs != "" {
		azsFilters = strings.Split(*filterAZs, ",")
	}
	azsFilter := filters.NewAZsFilter(azsFilters)

	var collectorsFilters []string
	if *filterCollectors != "" {
		collectorsFilters = strings.Split(*filterCollectors, ",")
	}
	collectorsFilter, err := filters.NewCollectorsFilter(collectorsFilters)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	var processesFilters []string
	if *sdProcessesRegexp != "" {
		processesFilters = []string{*sdProcessesRegexp}
	}
	processesFilter, err := filters.NewRegexpFilter(processesFilters)
	if err != nil {
		log.Errorf("Error processing Processes Regexp: %v", err)
		os.Exit(1)
	}

	boshCollector := collectors.NewBoshCollector(
		*metricsNamespace,
		*sdFilename,
		deploymentsFetcher,
		collectorsFilter,
		azsFilter,
		processesFilter,
	)
	prometheus.MustRegister(boshCollector)

	http.Handle(*metricsPath, prometheus.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
             <head><title>BOSH Exporter</title></head>
             <body>
             <h1>BOSH Exporter</h1>
             <p><a href='` + *metricsPath + `'>Metrics</a></p>
             </body>
             </html>`))
	})

	log.Infoln("Listening on", *listenAddress)
	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}
