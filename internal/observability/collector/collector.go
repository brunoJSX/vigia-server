// Package collector holds the Collector mechanism — part of the
// processing_pipeline (see domain-discovery.md / spec.yaml). It is an
// internal, provisional label, not a domain concept: it obtains the raw data
// of a verification (operational data), it does not interpret it.
package collector

import (
	"context"
	"net/http"
	"time"
)

// Sample is raw operational data produced by a single verification — a
// "dado coletado" (samples / resultados de verificação / métricas). It has no
// identity of its own; it is evidence consumed by the Analyzer.
type Sample struct {
	Timestamp time.Time
	Success   bool
	Latency   time.Duration
}

// Collector obtains the raw data of a verification for a given target,
// using the configuration (target, timeout) provided by a Monitor.
type Collector interface {
	Collect(ctx context.Context, target string, timeout time.Duration) (Sample, error)
}

// HTTPCollector probes a target over HTTP — enough to cover Uptime and
// Checkout Monitors (RN-024); Dependency Monitors use the same interface.
type HTTPCollector struct {
	client *http.Client
}

func NewHTTPCollector(client *http.Client) *HTTPCollector {
	if client == nil {
		client = http.DefaultClient
	}
	return &HTTPCollector{client: client}
}

func (c *HTTPCollector) Collect(ctx context.Context, target string, timeout time.Duration) (Sample, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, target, nil)
	if err != nil {
		return Sample{}, err
	}

	start := time.Now()
	resp, err := c.client.Do(req)
	latency := time.Since(start)
	if err != nil {
		return Sample{Timestamp: start, Success: false, Latency: latency}, nil
	}
	defer resp.Body.Close()

	success := resp.StatusCode >= http.StatusOK && resp.StatusCode < http.StatusBadRequest
	return Sample{Timestamp: start, Success: success, Latency: latency}, nil
}
