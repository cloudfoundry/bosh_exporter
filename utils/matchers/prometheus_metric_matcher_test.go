package matchers_test

import (
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/cloudfoundry/bosh_exporter/utils/matchers"
)

var _ = ginkgo.Describe("matchers.PrometheusMetric", func() {
	var (
		metricNamespace       = "fake_namespace"
		metricSubsystem       = "fake_sybsystem"
		metricName            = "fake_name"
		metricHelp            = "Fake Metric Help"
		metricLabelName       = "fake_label_name"
		metricLabelValue      = "fake_label_value"
		metricConstLabelName  = "fake_constant_label_name"
		metricConstLabelValue = "fake_constant_label_value"
	)

	ginkgo.Context("When asserting equality between Counter Metrics", func() {
		ginkgo.It("should do the right thing", func() {
			expectedMetric := prometheus.NewCounter(
				prometheus.CounterOpts{
					Namespace:   metricNamespace,
					Subsystem:   metricSubsystem,
					Name:        metricName,
					Help:        metricHelp,
					ConstLabels: prometheus.Labels{metricConstLabelName: metricConstLabelValue},
				})
			expectedMetric.Inc()

			actualMetric := prometheus.NewCounter(
				prometheus.CounterOpts{
					Namespace:   metricNamespace,
					Subsystem:   metricSubsystem,
					Name:        metricName,
					Help:        metricHelp,
					ConstLabels: prometheus.Labels{metricConstLabelName: metricConstLabelValue},
				})
			actualMetric.Inc()

			// Ignore creation timestamps
			matchers.ResetMetricCounterCreationTimestamp(expectedMetric)
			matchers.ResetMetricCounterCreationTimestamp(actualMetric)

			gomega.Expect(expectedMetric).To(matchers.PrometheusMetric(actualMetric))
		})
	})

	ginkgo.Context("When asserting equality between CounterVec Metrics", func() {
		ginkgo.It("should do the right thing", func() {
			expectedMetric := prometheus.NewCounterVec(
				prometheus.CounterOpts{
					Namespace:   metricNamespace,
					Subsystem:   metricSubsystem,
					Name:        metricName,
					Help:        metricHelp,
					ConstLabels: prometheus.Labels{metricConstLabelName: metricConstLabelValue},
				},
				[]string{metricLabelName},
			)
			expectedMetric.WithLabelValues(metricLabelValue).Inc()

			actualMetric := prometheus.NewCounterVec(
				prometheus.CounterOpts{
					Namespace:   metricNamespace,
					Subsystem:   metricSubsystem,
					Name:        metricName,
					Help:        metricHelp,
					ConstLabels: prometheus.Labels{metricConstLabelName: metricConstLabelValue},
				},
				[]string{metricLabelName},
			)
			actualMetric.WithLabelValues(metricLabelValue).Inc()

			// Ignore creation timestamps
			matchers.ResetMetricCounterCreationTimestamp(expectedMetric.WithLabelValues(metricLabelValue))
			matchers.ResetMetricCounterCreationTimestamp(actualMetric.WithLabelValues(metricLabelValue))

			gomega.Expect(expectedMetric.WithLabelValues(metricLabelValue)).To(matchers.PrometheusMetric(actualMetric.WithLabelValues(metricLabelValue)))
		})
	})

	ginkgo.Context("When asserting equality between Gauge Metrics", func() {
		ginkgo.It("should do the right thing", func() {
			expectedMetric := prometheus.NewGauge(
				prometheus.GaugeOpts{
					Namespace:   metricNamespace,
					Subsystem:   metricSubsystem,
					Name:        metricName,
					Help:        metricHelp,
					ConstLabels: prometheus.Labels{metricConstLabelName: metricConstLabelValue},
				})
			expectedMetric.Inc()

			actualMetric := prometheus.NewGauge(
				prometheus.GaugeOpts{
					Namespace:   metricNamespace,
					Subsystem:   metricSubsystem,
					Name:        metricName,
					Help:        metricHelp,
					ConstLabels: prometheus.Labels{metricConstLabelName: metricConstLabelValue},
				})
			actualMetric.Inc()

			gomega.Expect(expectedMetric).To(matchers.PrometheusMetric(actualMetric))
		})
	})

	ginkgo.Context("When asserting equality between GaugeVec Metrics", func() {
		ginkgo.It("should do the right thing", func() {
			expectedMetric := prometheus.NewGaugeVec(
				prometheus.GaugeOpts{
					Namespace:   metricNamespace,
					Subsystem:   metricSubsystem,
					Name:        metricName,
					Help:        metricHelp,
					ConstLabels: prometheus.Labels{metricConstLabelName: metricConstLabelValue},
				},
				[]string{metricLabelName},
			)
			expectedMetric.WithLabelValues(metricLabelValue).Inc()

			actualMetric := prometheus.NewGaugeVec(
				prometheus.GaugeOpts{
					Namespace:   metricNamespace,
					Subsystem:   metricSubsystem,
					Name:        metricName,
					Help:        metricHelp,
					ConstLabels: prometheus.Labels{metricConstLabelName: metricConstLabelValue},
				},
				[]string{metricLabelName},
			)
			actualMetric.WithLabelValues(metricLabelValue).Inc()

			gomega.Expect(expectedMetric.WithLabelValues(metricLabelValue)).To(matchers.PrometheusMetric(actualMetric.WithLabelValues(metricLabelValue)))
		})
	})

	ginkgo.Context("When asserting equality between Histogram Metrics", func() {
		ginkgo.It("should do the right thing", func() {
			expectedMetric := prometheus.NewHistogram(
				prometheus.HistogramOpts{
					Namespace:   metricNamespace,
					Subsystem:   metricSubsystem,
					Name:        metricName,
					Help:        metricHelp,
					ConstLabels: prometheus.Labels{metricConstLabelName: metricConstLabelValue},
				})
			expectedMetric.Observe(float64(1))

			actualMetric := prometheus.NewHistogram(
				prometheus.HistogramOpts{
					Namespace:   metricNamespace,
					Subsystem:   metricSubsystem,
					Name:        metricName,
					Help:        metricHelp,
					ConstLabels: prometheus.Labels{metricConstLabelName: metricConstLabelValue},
				})
			actualMetric.Observe(float64(1))

			// Ignore creation timestamps
			matchers.ResetMetricHistogramCreationTimestamp(expectedMetric)
			matchers.ResetMetricHistogramCreationTimestamp(actualMetric)

			gomega.Expect(expectedMetric).To(matchers.PrometheusMetric(actualMetric))
		})
	})

	ginkgo.Context("When asserting equality between HistogramVec Metrics", func() {
		ginkgo.It("should do the right thing", func() {
			expectedMetric := prometheus.NewHistogramVec(
				prometheus.HistogramOpts{
					Namespace:   metricNamespace,
					Subsystem:   metricSubsystem,
					Name:        metricName,
					Help:        metricHelp,
					ConstLabels: prometheus.Labels{metricConstLabelName: metricConstLabelValue},
				},
				[]string{metricLabelName},
			)
			expectedMetric.WithLabelValues(metricLabelValue).Observe(float64(1))

			actualMetric := prometheus.NewHistogramVec(
				prometheus.HistogramOpts{
					Namespace:   metricNamespace,
					Subsystem:   metricSubsystem,
					Name:        metricName,
					Help:        metricHelp,
					ConstLabels: prometheus.Labels{metricConstLabelName: metricConstLabelValue},
				},
				[]string{metricLabelName},
			)
			actualMetric.WithLabelValues(metricLabelValue).Observe(float64(1))

			// Ignore creation timestamps
			matchers.ResetMetricHistogramCreationTimestamp(expectedMetric.WithLabelValues(metricLabelValue).(prometheus.Histogram))
			matchers.ResetMetricHistogramCreationTimestamp(actualMetric.WithLabelValues(metricLabelValue).(prometheus.Histogram))

			gomega.Expect(expectedMetric.WithLabelValues(metricLabelValue)).To(matchers.PrometheusMetric(actualMetric.WithLabelValues(metricLabelValue).(prometheus.Histogram)))
		})
	})

	ginkgo.Context("When asserting equality between Summary Metrics", func() {
		ginkgo.It("should do the right thing", func() {
			expectedMetric := prometheus.NewSummary(
				prometheus.SummaryOpts{
					Namespace:   metricNamespace,
					Subsystem:   metricSubsystem,
					Name:        metricName,
					Help:        metricHelp,
					ConstLabels: prometheus.Labels{metricConstLabelName: metricConstLabelValue},
				})
			expectedMetric.Observe(float64(1))

			actualMetric := prometheus.NewSummary(
				prometheus.SummaryOpts{
					Namespace:   metricNamespace,
					Subsystem:   metricSubsystem,
					Name:        metricName,
					Help:        metricHelp,
					ConstLabels: prometheus.Labels{metricConstLabelName: metricConstLabelValue},
				})
			actualMetric.Observe(float64(1))

			// Ignore creation timestamps
			matchers.ResetMetricSummaryCreationTimestamp(expectedMetric)
			matchers.ResetMetricSummaryCreationTimestamp(actualMetric)

			gomega.Expect(expectedMetric).To(matchers.PrometheusMetric(actualMetric))
		})
	})

	ginkgo.Context("When asserting equality between SummaryVec Metrics", func() {
		ginkgo.It("should do the right thing", func() {
			expectedMetric := prometheus.NewSummaryVec(
				prometheus.SummaryOpts{
					Namespace:   metricNamespace,
					Subsystem:   metricSubsystem,
					Name:        metricName,
					Help:        metricHelp,
					ConstLabels: prometheus.Labels{metricConstLabelName: metricConstLabelValue},
				},
				[]string{metricLabelName},
			)
			expectedMetric.WithLabelValues(metricLabelValue).Observe(float64(1))

			actualMetric := prometheus.NewSummaryVec(
				prometheus.SummaryOpts{
					Namespace:   metricNamespace,
					Subsystem:   metricSubsystem,
					Name:        metricName,
					Help:        metricHelp,
					ConstLabels: prometheus.Labels{metricConstLabelName: metricConstLabelValue},
				},
				[]string{metricLabelName},
			)
			actualMetric.WithLabelValues(metricLabelValue).Observe(float64(1))

			// Ignore creation timestamps
			matchers.ResetMetricSummaryCreationTimestamp(expectedMetric.WithLabelValues(metricLabelValue).(prometheus.Summary))
			matchers.ResetMetricSummaryCreationTimestamp(actualMetric.WithLabelValues(metricLabelValue).(prometheus.Summary))

			gomega.Expect(expectedMetric.WithLabelValues(metricLabelValue)).To(matchers.PrometheusMetric(actualMetric.WithLabelValues(metricLabelValue).(prometheus.Summary)))
		})
	})
})
