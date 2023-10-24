package zeebe

import (
	"github.com/camunda/zeebe/clients/go/v8/pkg/worker"
	"time"
)

type WorkerOption interface {
	apply(builder worker.JobWorkerBuilderStep3)
}

type ownerNameWorkerOption struct {
	name string
}

type timeoutOption struct {
	timeout time.Duration
}

type requestTimeoutOption struct {
	requestTimeout time.Duration
}

type maxJobsActiveOption struct {
	maxJobsActive int
}

type concurrencyOption struct {
	concurrency int
}

type pollIntervalOption struct {
	pollInterval time.Duration
}

type pollThresholdOption struct {
	pollThreshold float64
}

type fetchVariablesOption struct {
	fetchVariables []string
}

type metricsOption struct {
	metrics worker.JobWorkerMetrics
}

func (o *ownerNameWorkerOption) apply(builder worker.JobWorkerBuilderStep3) {
	builder.Name(o.name)
}

func (o *timeoutOption) apply(builder worker.JobWorkerBuilderStep3) {
	builder.Timeout(o.timeout)
}

func (o *requestTimeoutOption) apply(builder worker.JobWorkerBuilderStep3) {
	builder.RequestTimeout(o.requestTimeout)
}

func (o *maxJobsActiveOption) apply(builder worker.JobWorkerBuilderStep3) {
	builder.MaxJobsActive(o.maxJobsActive)
}

func (o *concurrencyOption) apply(builder worker.JobWorkerBuilderStep3) {
	builder.Concurrency(o.concurrency)
}

func (o *pollIntervalOption) apply(builder worker.JobWorkerBuilderStep3) {
	builder.PollInterval(o.pollInterval)
}

func (o *pollThresholdOption) apply(builder worker.JobWorkerBuilderStep3) {
	builder.PollThreshold(o.pollThreshold)
}

func (o *fetchVariablesOption) apply(builder worker.JobWorkerBuilderStep3) {
	builder.FetchVariables(o.fetchVariables...)
}

func (o *metricsOption) apply(builder worker.JobWorkerBuilderStep3) {
	builder.Metrics(o.metrics)
}

func WithOwnerName(ownerName string) WorkerOption {
	return &ownerNameWorkerOption{name: ownerName}
}

func WithTimeout(timeout time.Duration) WorkerOption {
	return &timeoutOption{timeout: timeout}
}

func WithRequestTimeout(requestTimeout time.Duration) WorkerOption {
	return &requestTimeoutOption{requestTimeout: requestTimeout}
}

func WithMaxJobsActive(maxJobsActive int) WorkerOption {
	return &maxJobsActiveOption{maxJobsActive: maxJobsActive}
}

func WithConcurrency(concurrency int) WorkerOption {
	return &concurrencyOption{concurrency: concurrency}
}

func WithPollInterval(pollInterval time.Duration) WorkerOption {
	return &pollIntervalOption{pollInterval: pollInterval}
}

func WithPollThreshold(pollThreshold float64) WorkerOption {
	return &pollThresholdOption{pollThreshold: pollThreshold}
}

func WithFetchVariables(fetchVariables []string) WorkerOption {
	return &fetchVariablesOption{fetchVariables: fetchVariables}
}

func WithMetrics(metrics worker.JobWorkerMetrics) WorkerOption {
	return &metricsOption{metrics: metrics}
}
