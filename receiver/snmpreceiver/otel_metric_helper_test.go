// Copyright  The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package snmpreceiver // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/snmpreceiver"

import (
	"testing"
	"time"

	// client is an autogenerated mock type for the client type
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/receiver"
)

func TestGetResourceKey(t *testing.T) {
	testCases := []struct {
		desc     string
		testFunc func(*testing.T)
	}{
		{
			desc: "Empty arguments gives empty key",
			testFunc: func(t *testing.T) {
				expectedKey := ""
				actualKey := getResourceKey([]string{}, "")
				require.Equal(t, expectedKey, actualKey)
			},
		},
		{
			desc: "Returns stringified key from slice and index",
			testFunc: func(t *testing.T) {
				expectedKey := "key1,key2.1"
				actualKey := getResourceKey([]string{"key1", "key2"}, ".1")
				require.Equal(t, expectedKey, actualKey)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, tc.testFunc)
	}
}

func TestNewOTELMetricHelper(t *testing.T) {
	testCases := []struct {
		desc     string
		testFunc func(*testing.T)
	}{
		{
			desc: "Returns a good otelMetricHelper",
			testFunc: func(t *testing.T) {
				settings := receiver.CreateSettings{}
				helper := newOTELMetricHelper(settings, pcommon.NewTimestampFromTime(time.Now()))
				require.NotNil(t, helper)
				require.NotNil(t, helper.metrics)
				require.NotNil(t, helper.resourceMetricsSlice)
				require.NotNil(t, helper.dataPointStartTime)
				require.NotNil(t, helper.dataPointTime)
				require.Equal(t, settings, helper.settings)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, tc.testFunc)
	}
}

func TestGetResource(t *testing.T) {
	testCases := []struct {
		desc     string
		testFunc func(*testing.T)
	}{
		{
			desc: "Returns nil when resource not yet created",
			testFunc: func(t *testing.T) {
				settings := receiver.CreateSettings{}
				helper := newOTELMetricHelper(settings, pcommon.NewTimestampFromTime(time.Now()))
				actual := helper.getResource("r1")
				require.Nil(t, actual)
			},
		},
		{
			desc: "Returns resource when already created",
			testFunc: func(t *testing.T) {
				settings := receiver.CreateSettings{}
				helper := newOTELMetricHelper(settings, pcommon.NewTimestampFromTime(time.Now()))
				resource := helper.resourceMetricsSlice.AppendEmpty()
				resource.Resource().Attributes().PutStr("key1", "val1")
				helper.resourcesByKey["r1"] = &resource
				actual := helper.getResource("r1")
				require.Equal(t, &resource, actual)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, tc.testFunc)
	}
}

func TestCreateResource(t *testing.T) {
	testCases := []struct {
		desc     string
		testFunc func(*testing.T)
	}{
		{
			desc: "Creates resource with given attributes and saves it for easy reference",
			testFunc: func(t *testing.T) {
				settings := receiver.CreateSettings{}
				helper := newOTELMetricHelper(settings, pcommon.NewTimestampFromTime(time.Now()))
				actual := helper.createResource("r1", map[string]string{"key1": "val1"})
				require.NotNil(t, actual)
				val, exists := actual.Resource().Attributes().Get("key1")
				require.Equal(t, true, exists)
				require.Equal(t, "val1", val.AsString())
				require.Equal(t, actual, helper.resourcesByKey["r1"])
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, tc.testFunc)
	}
}

func TestGetMetric(t *testing.T) {
	testCases := []struct {
		desc     string
		testFunc func(*testing.T)
	}{
		{
			desc: "Returns nil when resource not yet created",
			testFunc: func(t *testing.T) {
				settings := receiver.CreateSettings{}
				helper := newOTELMetricHelper(settings, pcommon.NewTimestampFromTime(time.Now()))
				actual := helper.getMetric("r1", "m1")
				require.Nil(t, actual)
			},
		},
		{
			desc: "Returns nil when metric not yet created",
			testFunc: func(t *testing.T) {
				settings := receiver.CreateSettings{}
				helper := newOTELMetricHelper(settings, pcommon.NewTimestampFromTime(time.Now()))
				resource := helper.resourceMetricsSlice.AppendEmpty()
				resource.Resource().Attributes().PutStr("key1", "val1")
				helper.resourcesByKey["r1"] = &resource
				helper.metricsByResource["r1"] = map[string]*pmetric.Metric{}
				actual := helper.getMetric("r1", "m1")
				require.Nil(t, actual)
			},
		},
		{
			desc: "Returns metric when already created",
			testFunc: func(t *testing.T) {
				settings := receiver.CreateSettings{}
				helper := newOTELMetricHelper(settings, pcommon.NewTimestampFromTime(time.Now()))
				resource := helper.resourceMetricsSlice.AppendEmpty()
				resource.Resource().Attributes().PutStr("key1", "val1")
				helper.resourcesByKey["r1"] = &resource
				metric := resource.ScopeMetrics().AppendEmpty().Metrics().AppendEmpty()
				metric.SetName("Metric 1")
				helper.metricsByResource["r1"] = map[string]*pmetric.Metric{}
				helper.metricsByResource["r1"]["m1"] = &metric
				actual := helper.getMetric("r1", "m1")
				require.Equal(t, &metric, actual)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, tc.testFunc)
	}
}

func TestCreateMetric(t *testing.T) {
	testCases := []struct {
		desc     string
		testFunc func(*testing.T)
	}{
		{
			desc: "Returns error when resource does not exist",
			testFunc: func(t *testing.T) {
				settings := receiver.CreateSettings{}
				helper := newOTELMetricHelper(settings, pcommon.NewTimestampFromTime(time.Now()))
				metricCfg := MetricConfig{
					Description: "description",
					Unit:        "1",
					Gauge: &GaugeMetric{
						ValueType: "int",
					},
				}
				actual, err := helper.createMetric("r1", "m1", &metricCfg)
				require.Nil(t, actual)
				require.EqualError(t, err, "cannot create metric 'm1' as no resource exists for it to be attached")
			},
		},
		{
			desc: "Creates gauge metric and saves it for easy reference",
			testFunc: func(t *testing.T) {
				settings := receiver.CreateSettings{}
				helper := newOTELMetricHelper(settings, pcommon.NewTimestampFromTime(time.Now()))
				resource := helper.resourceMetricsSlice.AppendEmpty()
				resource.ScopeMetrics().AppendEmpty()
				resource.Resource().Attributes().PutStr("key1", "val1")
				helper.resourcesByKey["r1"] = &resource
				helper.metricsByResource["r1"] = map[string]*pmetric.Metric{}
				metricCfg := MetricConfig{
					Description: "description",
					Unit:        "1",
					Gauge: &GaugeMetric{
						ValueType: "int",
					},
				}
				actual, err := helper.createMetric("r1", "m1", &metricCfg)
				require.NoError(t, err)
				require.NotNil(t, actual)
				require.Equal(t, "description", actual.Description())
				require.NotNil(t, actual.Gauge())
				require.Equal(t, "m1", actual.Name())
				require.Equal(t, "1", actual.Unit())
				require.Equal(t, actual, helper.metricsByResource["r1"]["m1"])
			},
		},
		{
			desc: "Creates sum metric and saves it for easy reference",
			testFunc: func(t *testing.T) {
				settings := receiver.CreateSettings{}
				helper := newOTELMetricHelper(settings, pcommon.NewTimestampFromTime(time.Now()))
				resource := helper.resourceMetricsSlice.AppendEmpty()
				resource.ScopeMetrics().AppendEmpty()
				resource.Resource().Attributes().PutStr("key1", "val1")
				helper.resourcesByKey["r1"] = &resource
				helper.metricsByResource["r1"] = map[string]*pmetric.Metric{}
				metricCfg := MetricConfig{
					Description: "description",
					Unit:        "1",
					Sum: &SumMetric{
						Aggregation: "delta",
						Monotonic:   false,
						ValueType:   "double",
					},
				}
				actual, err := helper.createMetric("r1", "m1", &metricCfg)
				require.NoError(t, err)
				require.NotNil(t, actual)
				require.Equal(t, "description", actual.Description())
				require.Equal(t, pmetric.AggregationTemporalityDelta, actual.Sum().AggregationTemporality())
				require.Equal(t, false, actual.Sum().IsMonotonic())
				require.Equal(t, "m1", actual.Name())
				require.Equal(t, "1", actual.Unit())
				require.Equal(t, actual, helper.metricsByResource["r1"]["m1"])
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, tc.testFunc)
	}
}

func TestAddMetricDataPoint(t *testing.T) {
	testCases := []struct {
		desc     string
		testFunc func(*testing.T)
	}{
		{
			desc: "Returns error when resource does not exist",
			testFunc: func(t *testing.T) {
				settings := receiver.CreateSettings{}
				helper := newOTELMetricHelper(settings, pcommon.NewTimestampFromTime(time.Now()))
				metricCfg := MetricConfig{
					Description: "description",
					Unit:        "1",
					Gauge: &GaugeMetric{
						ValueType: "int",
					},
				}
				data := SNMPData{
					valueType: integerVal,
					value:     int64(10),
				}
				attributes := map[string]string{"key1": "val1"}
				actual, err := helper.addMetricDataPoint("r2", "m2", &metricCfg, data, attributes)
				require.Nil(t, actual)
				require.EqualError(t, err, "cannot retrieve datapoints from metric 'm2' as it does not currently exist")
			},
		},
		{
			desc: "Returns error when metric does not exist",
			testFunc: func(t *testing.T) {
				settings := receiver.CreateSettings{}
				helper := newOTELMetricHelper(settings, pcommon.NewTimestampFromTime(time.Now()))
				resource := helper.resourceMetricsSlice.AppendEmpty()
				resource.ScopeMetrics().AppendEmpty()
				resource.Resource().Attributes().PutStr("key1", "val1")
				helper.resourcesByKey["r1"] = &resource
				helper.metricsByResource["r1"] = map[string]*pmetric.Metric{}
				metricCfg := MetricConfig{
					Description: "description",
					Unit:        "1",
					Gauge: &GaugeMetric{
						ValueType: "int",
					},
				}
				data := SNMPData{
					valueType: integerVal,
					value:     int64(10),
				}
				attributes := map[string]string{"key1": "val1"}
				actual, err := helper.addMetricDataPoint("r1", "m1", &metricCfg, data, attributes)
				require.Nil(t, actual)
				require.EqualError(t, err, "cannot retrieve datapoints from metric 'm1' as it does not currently exist")
			},
		},
		{
			desc: "Creates data points on existing gauge metric using passed in data",
			testFunc: func(t *testing.T) {
				settings := receiver.CreateSettings{}
				helper := newOTELMetricHelper(settings, pcommon.NewTimestampFromTime(time.Now()))
				resource := helper.resourceMetricsSlice.AppendEmpty()
				resource.ScopeMetrics().AppendEmpty()
				resource.Resource().Attributes().PutStr("key1", "val1")
				helper.resourcesByKey["r1"] = &resource
				metric := resource.ScopeMetrics().AppendEmpty().Metrics().AppendEmpty()
				metric.SetName("Metric 1")
				metric.SetEmptyGauge()
				helper.metricsByResource["r1"] = map[string]*pmetric.Metric{}
				helper.metricsByResource["r1"]["m1"] = &metric
				metricCfg := MetricConfig{
					Description: "description",
					Unit:        "1",
					Gauge: &GaugeMetric{
						ValueType: "int",
					},
				}
				data := SNMPData{
					valueType: integerVal,
					value:     int64(10),
				}
				attributes := map[string]string{"key1": "val1"}
				actual, err := helper.addMetricDataPoint("r1", "m1", &metricCfg, data, attributes)
				require.NoError(t, err)
				require.Equal(t, data.value, actual.IntValue())
				val, exists := actual.Attributes().Get("key1")
				require.Equal(t, true, exists)
				require.Equal(t, "val1", val.AsString())
				metricDataPoint := metric.Gauge().DataPoints().At(0)
				require.Equal(t, &metricDataPoint, actual)
			},
		},
		{
			desc: "Creates data points on existing sum metric using passed in data",
			testFunc: func(t *testing.T) {
				settings := receiver.CreateSettings{}
				helper := newOTELMetricHelper(settings, pcommon.NewTimestampFromTime(time.Now()))
				resource := helper.resourceMetricsSlice.AppendEmpty()
				resource.ScopeMetrics().AppendEmpty()
				resource.Resource().Attributes().PutStr("key1", "val1")
				helper.resourcesByKey["r1"] = &resource
				metric := resource.ScopeMetrics().AppendEmpty().Metrics().AppendEmpty()
				metric.SetName("Metric 1")
				metric.SetEmptySum()
				metric.Sum().SetAggregationTemporality(pmetric.AggregationTemporalityCumulative)
				metric.Sum().SetIsMonotonic(true)
				helper.metricsByResource["r1"] = map[string]*pmetric.Metric{}
				helper.metricsByResource["r1"]["m1"] = &metric
				metricCfg := MetricConfig{
					Description: "description",
					Unit:        "1",
					Sum: &SumMetric{
						Aggregation: "cumulative",
						Monotonic:   true,
						ValueType:   "double",
					},
				}
				data := SNMPData{
					valueType: floatVal,
					value:     float64(10.0),
				}
				attributes := map[string]string{"key1": "val1"}
				actual, err := helper.addMetricDataPoint("r1", "m1", &metricCfg, data, attributes)
				require.NoError(t, err)
				require.Equal(t, data.value, actual.DoubleValue())
				val, exists := actual.Attributes().Get("key1")
				require.Equal(t, true, exists)
				require.Equal(t, "val1", val.AsString())
				metricDataPoint := metric.Sum().DataPoints().At(0)
				require.Equal(t, &metricDataPoint, actual)
			},
		},
		{
			desc: "Creates data points on existing metric converting float to int",
			testFunc: func(t *testing.T) {
				settings := receiver.CreateSettings{}
				helper := newOTELMetricHelper(settings, pcommon.NewTimestampFromTime(time.Now()))
				resource := helper.resourceMetricsSlice.AppendEmpty()
				resource.ScopeMetrics().AppendEmpty()
				resource.Resource().Attributes().PutStr("key1", "val1")
				helper.resourcesByKey["r1"] = &resource
				metric := resource.ScopeMetrics().AppendEmpty().Metrics().AppendEmpty()
				metric.SetName("Metric 1")
				metric.SetEmptyGauge()
				helper.metricsByResource["r1"] = map[string]*pmetric.Metric{}
				helper.metricsByResource["r1"]["m1"] = &metric
				metricCfg := MetricConfig{
					Description: "description",
					Unit:        "1",
					Gauge: &GaugeMetric{
						ValueType: "int",
					},
				}
				data := SNMPData{
					valueType: floatVal,
					value:     float64(10.0),
				}
				attributes := map[string]string{"key1": "val1"}
				actual, err := helper.addMetricDataPoint("r1", "m1", &metricCfg, data, attributes)
				require.NoError(t, err)
				require.Equal(t, int64(10), actual.IntValue())
				val, exists := actual.Attributes().Get("key1")
				require.Equal(t, true, exists)
				require.Equal(t, "val1", val.AsString())
				metricDataPoint := metric.Gauge().DataPoints().At(0)
				require.Equal(t, &metricDataPoint, actual)
			},
		},
		{
			desc: "Creates data points on existing metric converting int to float",
			testFunc: func(t *testing.T) {
				settings := receiver.CreateSettings{}
				helper := newOTELMetricHelper(settings, pcommon.NewTimestampFromTime(time.Now()))
				resource := helper.resourceMetricsSlice.AppendEmpty()
				resource.ScopeMetrics().AppendEmpty()
				resource.Resource().Attributes().PutStr("key1", "val1")
				helper.resourcesByKey["r1"] = &resource
				metric := resource.ScopeMetrics().AppendEmpty().Metrics().AppendEmpty()
				metric.SetName("Metric 1")
				metric.SetEmptyGauge()
				helper.metricsByResource["r1"] = map[string]*pmetric.Metric{}
				helper.metricsByResource["r1"]["m1"] = &metric
				metricCfg := MetricConfig{
					Description: "description",
					Unit:        "1",
					Gauge: &GaugeMetric{
						ValueType: "double",
					},
				}
				data := SNMPData{
					valueType: integerVal,
					value:     int64(10),
				}
				attributes := map[string]string{"key1": "val1"}
				actual, err := helper.addMetricDataPoint("r1", "m1", &metricCfg, data, attributes)
				require.NoError(t, err)
				require.Equal(t, float64(10.0), actual.DoubleValue())
				val, exists := actual.Attributes().Get("key1")
				require.Equal(t, true, exists)
				require.Equal(t, "val1", val.AsString())
				metricDataPoint := metric.Gauge().DataPoints().At(0)
				require.Equal(t, &metricDataPoint, actual)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, tc.testFunc)
	}
}
