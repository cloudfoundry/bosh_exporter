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

func ResetMetricCounterCreationTimestamp(counter prometheus.Metric) {
	metric := &dto.Metric{}
	_ = counter.Write(metric)

	metric.Counter.CreatedTimestamp.Reset()
}

func ResetMetricHistogramCreationTimestamp(histogram prometheus.Metric) {
	metric := &dto.Metric{}
	_ = histogram.Write(metric)

	//	fmt.Printf("Before_Stamp: %v", metric.Histogram.CreatedTimestamp)
	metric.Histogram.CreatedTimestamp.Reset()
	// fmt.Printf("After_stamp: %v", metric.Histogram.CreatedTimestamp)
}

func ResetMetricSummaryCreationTimestamp(summary prometheus.Metric) {
	metric := &dto.Metric{}
	_ = summary.Write(metric)

	metric.Summary.CreatedTimestamp.Reset()
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
