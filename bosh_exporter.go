package main

import (
	"errors"
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

	"github.com/bosh-prometheus/bosh_exporter/collectors"
	"github.com/bosh-prometheus/bosh_exporter/deployments"
	"github.com/bosh-prometheus/bosh_exporter/filters"
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

	metricsEnvironment = flag.String(
		"metrics.environment", "",
		"Environment label to be attached to metrics ($BOSH_EXPORTER_METRICS_ENVIRONMENT).",
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

	authUsername = flag.String(
		"web.auth.username", "",
		"Username for web interface basic auth ($BOSH_EXPORTER_WEB_AUTH_USERNAME).",
	)

	authPassword = flag.String(
		"web.auth.password", "",
		"Password for web interface basic auth ($BOSH_EXPORTER_WEB_AUTH_PASSWORD).",
	)

	tlsCertFile = flag.String(
		"web.tls.cert_file", "",
		"Path to a file that contains the TLS certificate (PEM format). If the certificate is signed by a certificate authority, the file should be the concatenation of the server's certificate, any intermediates, and the CA's certificate ($BOSH_EXPORTER_WEB_TLS_CERTFILE).",
	)

	tlsKeyFile = flag.String(
		"web.tls.key_file", "",
		"Path to a file that contains the TLS private key (PEM format) ($BOSH_EXPORTER_WEB_TLS_KEYFILE).",
	)
)

func init() {
	prometheus.MustRegister(version.NewCollector(*metricsNamespace))
}

func overrideFlagsWithEnvVars() {
	overrideWithEnvVar("BOSH_EXPORTER_BOSH_URL", boshURL)
	overrideWithEnvVar("BOSH_EXPORTER_BOSH_USERNAME", boshUsername)
	overrideWithEnvVar("BOSH_EXPORTER_BOSH_PASSWORD", boshPassword)
	overrideWithEnvVar("BOSH_EXPORTER_BOSH_UAA_CLIENT_ID", boshUAAClientID)
	overrideWithEnvVar("BOSH_EXPORTER_BOSH_UAA_CLIENT_SECRET", boshUAAClientSecret)
	overrideWithEnvVar("BOSH_EXPORTER_BOSH_LOG_LEVEL", boshLogLevel)
	overrideWithEnvVar("BOSH_EXPORTER_BOSH_CA_CERT_FILE", boshCACertFile)
	overrideWithEnvVar("BOSH_EXPORTER_FILTER_DEPLOYMENTS", filterDeployments)
	overrideWithEnvVar("BOSH_EXPORTER_FILTER_AZS", filterAZs)
	overrideWithEnvVar("BOSH_EXPORTER_FILTER_COLLECTORS", filterCollectors)
	overrideWithEnvVar("BOSH_EXPORTER_METRICS_NAMESPACE", metricsNamespace)
	overrideWithEnvVar("BOSH_EXPORTER_METRICS_ENVIRONMENT", metricsEnvironment)
	overrideWithEnvVar("BOSH_EXPORTER_SD_FILENAME", sdFilename)
	overrideWithEnvVar("BOSH_EXPORTER_SD_PROCESSES_REGEXP", sdProcessesRegexp)
	overrideWithEnvVar("BOSH_EXPORTER_WEB_LISTEN_ADDRESS", listenAddress)
	overrideWithEnvVar("BOSH_EXPORTER_WEB_TELEMETRY_PATH", metricsPath)
	overrideWithEnvVar("BOSH_EXPORTER_WEB_AUTH_USERNAME", authUsername)
	overrideWithEnvVar("BOSH_EXPORTER_WEB_AUTH_PASSWORD", authPassword)
	overrideWithEnvVar("BOSH_EXPORTER_WEB_TLS_CERTFILE", tlsCertFile)
	overrideWithEnvVar("BOSH_EXPORTER_WEB_TLS_KEYFILE", tlsKeyFile)
}

func overrideWithEnvVar(name string, value *string) {
	envValue := os.Getenv(name)
	if envValue != "" {
		*value = envValue
	}
}

type basicAuthHandler struct {
	handler  http.HandlerFunc
	username string
	password string
}

func (h *basicAuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	username, password, ok := r.BasicAuth()
	if !ok || username != h.username || password != h.password {
		log.Errorf("Invalid HTTP auth from `%s`", r.RemoteAddr)
		w.Header().Set("WWW-Authenticate", "Basic realm=\"metrics\"")
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}
	h.handler(w, r)
	return
}

func prometheusHandler() http.Handler {
	handler := prometheus.Handler()

	if *authUsername != "" && *authPassword != "" {
		handler = &basicAuthHandler{
			handler:  prometheus.Handler().ServeHTTP,
			username: *authUsername,
			password: *authPassword,
		}
	}

	return handler
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

	anonymousDirector, err := director.NewFactory(logger).New(directorConfig, nil, nil)
	if err != nil {
		return nil, err
	}

	boshInfo, err := anonymousDirector.Info()
	if err != nil {
		return nil, err
	}

	if boshInfo.Auth.Type != "uaa" {
		directorConfig.Client = *boshUsername
		directorConfig.ClientSecret = *boshPassword
	} else {
		uaaURL := boshInfo.Auth.Options["url"]
		uaaURLStr, ok := uaaURL.(string)
		if !ok {
			return nil, errors.New(fmt.Sprintf("Expected UAA URL '%s' to be a string", uaaURL))
		}

		uaaConfig, err := uaa.NewConfigFromURL(uaaURLStr)
		if err != nil {
			return nil, err
		}

		uaaConfig.CACert = boshCACert

		if *boshUAAClientID != "" && *boshUAAClientSecret != "" {
			uaaConfig.Client = *boshUAAClientID
			uaaConfig.ClientSecret = *boshUAAClientSecret
		} else {
			uaaConfig.Client = "bosh_cli"
		}

		uaaFactory := uaa.NewFactory(logger)
		uaaClient, err := uaaFactory.New(uaaConfig)
		if err != nil {
			return nil, err
		}

		if *boshUAAClientID != "" && *boshUAAClientSecret != "" {
			directorConfig.TokenFunc = uaa.NewClientTokenSession(uaaClient).TokenFunc
		} else {
			answers := []uaa.PromptAnswer{
				uaa.PromptAnswer{
					Key:   "username",
					Value: *boshUsername,
				},
				uaa.PromptAnswer{
					Key:   "password",
					Value: *boshPassword,
				},
			}
			accessToken, err := uaaClient.OwnerPasswordCredentialsGrant(answers)
			if err != nil {
				return nil, err
			}

			origToken := uaaClient.NewStaleAccessToken(accessToken.RefreshToken().Value())
			directorConfig.TokenFunc = uaa.NewAccessTokenSession(origToken).TokenFunc
		}
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
		*metricsEnvironment,
		boshInfo.Name,
		boshInfo.UUID,
		*sdFilename,
		deploymentsFetcher,
		collectorsFilter,
		azsFilter,
		processesFilter,
	)
	prometheus.MustRegister(boshCollector)

	http.Handle(*metricsPath, prometheusHandler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
             <head><title>BOSH Exporter</title></head>
             <body>
             <h1>BOSH Exporter</h1>
             <p><a href='` + *metricsPath + `'>Metrics</a></p>
             </body>
             </html>`))
	})

	if *tlsCertFile != "" && *tlsKeyFile != "" {
		log.Infoln("Listening TLS on", *listenAddress)
		log.Fatal(http.ListenAndServeTLS(*listenAddress, *tlsCertFile, *tlsKeyFile, nil))
	} else {
		log.Infoln("Listening on", *listenAddress)
		log.Fatal(http.ListenAndServe(*listenAddress, nil))
	}
}
