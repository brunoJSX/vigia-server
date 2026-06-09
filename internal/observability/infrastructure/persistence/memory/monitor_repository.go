// Package memory provides in-memory MonitorRepository / IncidentRepository /
// SampleRepository implementations — used by application tests, and as a
// lightweight alternative to the Postgres persistence for local runs.
package memory

import (
	"context"
	"fmt"
	"sync"

	"github.com/vigia/vigia-v1/internal/observability/monitor"
)

type MonitorRepository struct {
	mu       sync.Mutex
	monitors map[string]monitor.Monitor
}

func NewMonitorRepository() *MonitorRepository {
	return &MonitorRepository{monitors: make(map[string]monitor.Monitor)}
}

func (r *MonitorRepository) Save(ctx context.Context, m monitor.Monitor) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.monitors[m.ID] = m
	return nil
}

func (r *MonitorRepository) FindByID(ctx context.Context, id string) (monitor.Monitor, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	m, ok := r.monitors[id]
	if !ok {
		return monitor.Monitor{}, fmt.Errorf("monitor %q not found", id)
	}
	return m, nil
}

func (r *MonitorRepository) FindActive(ctx context.Context) ([]monitor.Monitor, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var active []monitor.Monitor
	for _, m := range r.monitors {
		if m.IsActive() {
			active = append(active, m)
		}
	}
	return active, nil
}

func (r *MonitorRepository) FindAll(ctx context.Context) ([]monitor.Monitor, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	out := make([]monitor.Monitor, 0, len(r.monitors))
	for _, m := range r.monitors {
		out = append(out, m)
	}
	return out, nil
}
