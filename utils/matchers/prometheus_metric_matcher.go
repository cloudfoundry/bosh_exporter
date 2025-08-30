package matchers

import (
	"fmt"
	"reflect"

	"github.com/onsi/gomega/format"
	"github.com/onsi/gomega/types"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
)

func PrometheusMetric(expected prometheus.Metric) types.GomegaMatcher {
	expectedMetric := &dto.Metric{}
	_ = expected.Write(expectedMetric)

	return &PrometheusMetricMatcher{
		Desc:   expected.Desc(),
		Metric: expectedMetric,
	}
}

type PrometheusMetricMatcher struct {
	Desc   *prometheus.Desc
	Metric *dto.Metric
}

func (matcher *PrometheusMetricMatcher) Match(actual interface{}) (success bool, err error) {
	metric, ok := actual.(prometheus.Metric)
	if !ok {
		return false, fmt.Errorf("PrometheusMetric matcher expects a prometheus.Metric")
	}

	actualMetric := &dto.Metric{}
	_ = metric.Write(actualMetric)

	if actualMetric.Counter != nil && matcher.Metric.Counter != nil {
		actualMetric.Counter.CreatedTimestamp.Reset()
		matcher.Metric.Counter.CreatedTimestamp.Reset()
	}

	if actualMetric.Histogram != nil && matcher.Metric.Histogram != nil {
		actualMetric.Histogram.CreatedTimestamp.Reset()
		matcher.Metric.Histogram.CreatedTimestamp.Reset()
	}

	if actualMetric.Summary != nil && matcher.Metric.Summary != nil {
		actualMetric.Summary.CreatedTimestamp.Reset()
		matcher.Metric.Summary.CreatedTimestamp.Reset()
	}

	if !reflect.DeepEqual(metric.Desc().String(), matcher.Desc.String()) {
		return false, nil
	}

	return reflect.DeepEqual(actualMetric.String(), matcher.Metric.String()), nil
}

func (matcher *PrometheusMetricMatcher) FailureMessage(actual interface{}) (message string) {
	metric, ok := actual.(prometheus.Metric)
	if ok {
		actualMetric := &dto.Metric{}
		_ = metric.Write(actualMetric)
		return format.Message(
			fmt.Sprintf("\n%s\nMetric{%s}", metric.Desc().String(), actualMetric.String()),
			"to equal",
			fmt.Sprintf("\n%s\nMetric{%s}", matcher.Desc.String(), matcher.Metric.String()),
		)
	}

	return format.Message(actual, "to equal", matcher)
}

func (matcher *PrometheusMetricMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return format.Message(actual, "not to equal", matcher)
}
