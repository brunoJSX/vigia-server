package application_test

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/vigia/vigia-v1/internal/observability/application"
	"github.com/vigia/vigia-v1/internal/observability/collector"
	"github.com/vigia/vigia-v1/internal/shared/clock"
	"github.com/vigia/vigia-v1/internal/shared/id"
)

// fakeCollector returns a pre-scripted sequence of success/failure results,
// one per call — enough to drive the Analyzer's consecutiveness rule
// (RN-027) deterministically across several CheckMonitor.Execute calls.
type fakeCollector struct {
	results []bool
	start   time.Time
	calls   int
}

func (c *fakeCollector) Collect(ctx context.Context, target string, timeout time.Duration) (collector.Sample, error) {
	success := true
	if c.calls < len(c.results) {
		success = c.results[c.calls]
	}
	ts := c.start.Add(time.Duration(c.calls) * time.Minute)
	c.calls++
	return collector.Sample{Timestamp: ts, Success: success, Latency: 10 * time.Millisecond}, nil
}

// spyPublisher records every Event it receives, for assertions on what the
// pipeline emitted.
type spyPublisher struct {
	mu     sync.Mutex
	events []application.Event
}

func (p *spyPublisher) Publish(ctx context.Context, e application.Event) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.events = append(p.events, e)
	return nil
}

func fixedClock(t time.Time) clock.Clock {
	return func() time.Time { return t }
}

func sequentialIDs(prefix string) id.Generator {
	n := 0
	return func() string {
		n++
		return fmt.Sprintf("%s-%d", prefix, n)
	}
}
