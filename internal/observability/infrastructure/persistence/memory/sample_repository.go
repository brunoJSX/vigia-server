package memory

import (
	"context"
	"sync"
	"time"

	"github.com/vigia/vigia-v1/internal/observability/collector"
)

type SampleRepository struct {
	mu      sync.Mutex
	samples map[string][]collector.Sample
}

func NewSampleRepository() *SampleRepository {
	return &SampleRepository{samples: make(map[string][]collector.Sample)}
}

func (r *SampleRepository) Save(ctx context.Context, monitorID string, s collector.Sample) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.samples[monitorID] = append(r.samples[monitorID], s)
	return nil
}

// FindRecent returns up to `limit` of the most recent samples, oldest first —
// assumes Save is called in chronological order, as CheckMonitor does.
func (r *SampleRepository) FindRecent(ctx context.Context, monitorID string, limit int) ([]collector.Sample, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	all := r.samples[monitorID]
	if limit <= 0 || limit > len(all) {
		limit = len(all)
	}

	out := make([]collector.Sample, limit)
	copy(out, all[len(all)-limit:])
	return out, nil
}

func (r *SampleRepository) FindByMonitorAndPeriod(ctx context.Context, monitorID string, from, to time.Time) ([]collector.Sample, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var out []collector.Sample
	for _, s := range r.samples[monitorID] {
		if !s.Timestamp.Before(from) && s.Timestamp.Before(to) {
			out = append(out, s)
		}
	}
	return out, nil
}
