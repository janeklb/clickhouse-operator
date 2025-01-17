// Copyright 2019 Altinity Ltd and/or its affiliates. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package metrics

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"

	"github.com/altinity/clickhouse-operator/pkg/metrics/operator"
)

// Metrics is a set of metrics that are tracked by the operator
type Metrics struct {
	// CHIReconcilesStarted is a number (counter) of started CHI reconciles
	CHIReconcilesStarted metric.Int64Counter
	// CHIReconcilesCompleted is a number (counter) of completed CHI reconciles.
	// In ideal world number of completed reconciles should be equal to CHIReconcilesStarted
	CHIReconcilesCompleted metric.Int64Counter
	// CHIReconcilesAborted is a number (counter) of explicitly aborted CHI reconciles.
	// This counter does not includes reconciles that we not completed due to external rasons, such as operator restart
	CHIReconcilesAborted metric.Int64Counter
	// CHIReconcilesTimings is a histogram of durations of successfully completed CHI reconciles
	CHIReconcilesTimings metric.Float64Histogram

	// HostReconcilesStarted is a number (counter) of started host reconciles
	HostReconcilesStarted metric.Int64Counter
	// HostReconcilesCompleted is a number (counter) of completed host reconciles.
	// In ideal world number of completed reconciles should be equal to HostReconcilesStarted
	HostReconcilesCompleted metric.Int64Counter
	// HostReconcilesRestarts is a number (counter) of host restarts during reconcile
	HostReconcilesRestarts metric.Int64Counter
	// HostReconcilesErrors is a number (counter) of failed (non-completed) host reconciles.
	HostReconcilesErrors metric.Int64Counter
	// HostReconcilesTimings is a histogram of durations of successfully completed host reconciles
	HostReconcilesTimings metric.Float64Histogram

	PodAddEvents    metric.Int64Counter
	PodUpdateEvents metric.Int64Counter
	PodDeleteEvents metric.Int64Counter
}

var m *Metrics

func createMetrics() *Metrics {
	// The unit u should be defined using the appropriate [UCUM](https://ucum.org) case-sensitive code.
	CHIReconcilesStarted, _ := operator.Meter().Int64Counter(
		"clickhouse_operator_chi_reconciles_started",
		metric.WithDescription("number of CHI reconciles started"),
		metric.WithUnit("items"),
	)
	CHIReconcilesCompleted, _ := operator.Meter().Int64Counter(
		"clickhouse_operator_chi_reconciles_completed",
		metric.WithDescription("number of CHI reconciles completed successfully"),
		metric.WithUnit("items"),
	)
	CHIReconcilesAborted, _ := operator.Meter().Int64Counter(
		"clickhouse_operator_chi_reconciles_aborted",
		metric.WithDescription("number of CHI reconciles aborted"),
		metric.WithUnit("items"),
	)
	CHIReconcilesTimings, _ := operator.Meter().Float64Histogram(
		"clickhouse_operator_chi_reconciles_timings",
		metric.WithDescription("timings of CHI reconciles completed successfully"),
		metric.WithUnit("s"),
	)

	HostReconcilesStarted, _ := operator.Meter().Int64Counter(
		"clickhouse_operator_host_reconciles_started",
		metric.WithDescription("number of host reconciles started"),
		metric.WithUnit("items"),
	)
	HostReconcilesCompleted, _ := operator.Meter().Int64Counter(
		"clickhouse_operator_host_reconciles_completed",
		metric.WithDescription("number of host reconciles completed successfully"),
		metric.WithUnit("items"),
	)
	HostReconcilesRestarts, _ := operator.Meter().Int64Counter(
		"clickhouse_operator_host_reconciles_restarts",
		metric.WithDescription("number of host restarts during reconciles"),
		metric.WithUnit("items"),
	)
	HostReconcilesErrors, _ := operator.Meter().Int64Counter(
		"clickhouse_operator_host_reconciles_errors",
		metric.WithDescription("number of host reconciles errors"),
		metric.WithUnit("items"),
	)
	HostReconcilesTimings, _ := operator.Meter().Float64Histogram(
		"clickhouse_operator_host_reconciles_timings",
		metric.WithDescription("timings of host reconciles completed successfully"),
		metric.WithUnit("s"),
	)

	PodAddEvents, _ := operator.Meter().Int64Counter(
		"clickhouse_operator_pod_add_events",
		metric.WithDescription("number PodAdd events"),
		metric.WithUnit("items"),
	)
	PodUpdateEvents, _ := operator.Meter().Int64Counter(
		"clickhouse_operator_pod_update_events",
		metric.WithDescription("number PodUpdate events"),
		metric.WithUnit("items"),
	)
	PodDeleteEvents, _ := operator.Meter().Int64Counter(
		"clickhouse_operator_pod_delete_events",
		metric.WithDescription("number PodDelete events"),
		metric.WithUnit("items"),
	)

	return &Metrics{
		CHIReconcilesStarted:   CHIReconcilesStarted,
		CHIReconcilesCompleted: CHIReconcilesCompleted,
		CHIReconcilesAborted:   CHIReconcilesAborted,
		CHIReconcilesTimings:   CHIReconcilesTimings,

		HostReconcilesStarted:   HostReconcilesStarted,
		HostReconcilesCompleted: HostReconcilesCompleted,
		HostReconcilesRestarts:  HostReconcilesRestarts,
		HostReconcilesErrors:    HostReconcilesErrors,
		HostReconcilesTimings:   HostReconcilesTimings,

		PodAddEvents:    PodAddEvents,
		PodUpdateEvents: PodUpdateEvents,
		PodDeleteEvents: PodDeleteEvents,
	}
}

func ensureMetrics() *Metrics {
	if m == nil {
		m = createMetrics()
	}
	return m
}

type BaseInfoGetter interface {
	GetName() string
	GetNamespace() string
	GetLabels() map[string]string
	GetAnnotations() map[string]string
}

func prepareLabels(cr BaseInfoGetter) (attributes []attribute.KeyValue) {
	labels, values := operator.GetMandatoryLabelsAndValues(cr)
	for i := range labels {
		label := labels[i]
		value := values[i]
		attributes = append(attributes, attribute.String(label, value))
	}

	return attributes
}

// metricsCHIInitZeroValues initializes all metrics for CHI to zero values if not already present with appropriate labels
//
// This is due to `rate` prometheus function limitation where it expects the metric to be 0-initialized with all possible labels
// and doesn't default to 0 if the metric is not present.
func CHIInitZeroValues(ctx context.Context, chi BaseInfoGetter) {
	ensureMetrics().CHIReconcilesStarted.Add(ctx, 0, metric.WithAttributes(prepareLabels(chi)...))
	ensureMetrics().CHIReconcilesCompleted.Add(ctx, 0, metric.WithAttributes(prepareLabels(chi)...))
	ensureMetrics().CHIReconcilesAborted.Add(ctx, 0, metric.WithAttributes(prepareLabels(chi)...))

	ensureMetrics().HostReconcilesStarted.Add(ctx, 0, metric.WithAttributes(prepareLabels(chi)...))
	ensureMetrics().HostReconcilesCompleted.Add(ctx, 0, metric.WithAttributes(prepareLabels(chi)...))
	ensureMetrics().HostReconcilesRestarts.Add(ctx, 0, metric.WithAttributes(prepareLabels(chi)...))
	ensureMetrics().HostReconcilesErrors.Add(ctx, 0, metric.WithAttributes(prepareLabels(chi)...))
}

func CHIReconcilesStarted(ctx context.Context, chi BaseInfoGetter) {
	ensureMetrics().CHIReconcilesStarted.Add(ctx, 1, metric.WithAttributes(prepareLabels(chi)...))
}
func CHIReconcilesCompleted(ctx context.Context, chi BaseInfoGetter) {
	ensureMetrics().CHIReconcilesCompleted.Add(ctx, 1, metric.WithAttributes(prepareLabels(chi)...))
}
func CHIReconcilesAborted(ctx context.Context, chi BaseInfoGetter) {
	ensureMetrics().CHIReconcilesAborted.Add(ctx, 1, metric.WithAttributes(prepareLabels(chi)...))
}
func CHIReconcilesTimings(ctx context.Context, chi BaseInfoGetter, seconds float64) {
	ensureMetrics().CHIReconcilesTimings.Record(ctx, seconds, metric.WithAttributes(prepareLabels(chi)...))
}

func HostReconcilesStarted(ctx context.Context, chi BaseInfoGetter) {
	ensureMetrics().HostReconcilesStarted.Add(ctx, 1, metric.WithAttributes(prepareLabels(chi)...))
}
func HostReconcilesCompleted(ctx context.Context, chi BaseInfoGetter) {
	ensureMetrics().HostReconcilesCompleted.Add(ctx, 1, metric.WithAttributes(prepareLabels(chi)...))
}
func HostReconcilesRestart(ctx context.Context, chi BaseInfoGetter) {
	ensureMetrics().HostReconcilesRestarts.Add(ctx, 1, metric.WithAttributes(prepareLabels(chi)...))
}
func HostReconcilesErrors(ctx context.Context, chi BaseInfoGetter) {
	ensureMetrics().HostReconcilesErrors.Add(ctx, 1, metric.WithAttributes(prepareLabels(chi)...))
}
func HostReconcilesTimings(ctx context.Context, chi BaseInfoGetter, seconds float64) {
	ensureMetrics().HostReconcilesTimings.Record(ctx, seconds, metric.WithAttributes(prepareLabels(chi)...))
}

func PodAdd(ctx context.Context) {
	ensureMetrics().PodAddEvents.Add(ctx, 1)
}
func metricsPodUpdate(ctx context.Context) {
	ensureMetrics().PodUpdateEvents.Add(ctx, 1)
}
func PodDelete(ctx context.Context) {
	ensureMetrics().PodDeleteEvents.Add(ctx, 1)
}
