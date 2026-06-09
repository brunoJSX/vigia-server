package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/vigia/vigia-v1/internal/observability/collector"
)

type SampleRepository struct {
	pool *pgxpool.Pool
}

func NewSampleRepository(pool *pgxpool.Pool) *SampleRepository {
	return &SampleRepository{pool: pool}
}

func (r *SampleRepository) Save(ctx context.Context, monitorID string, s collector.Sample) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO samples (monitor_id, "timestamp", success, latency_ns)
		VALUES ($1, $2, $3, $4)
	`, monitorID, s.Timestamp, s.Success, int64(s.Latency))
	return err
}

// FindRecent returns up to `limit` of the most recent samples, oldest first —
// matching the in-memory repository's contract that CheckMonitor relies on.
func (r *SampleRepository) FindRecent(ctx context.Context, monitorID string, limit int) ([]collector.Sample, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT "timestamp", success, latency_ns
		FROM samples
		WHERE monitor_id = $1
		ORDER BY "timestamp" DESC
		LIMIT $2
	`, monitorID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []collector.Sample
	for rows.Next() {
		s, err := scanSample(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	for i, j := 0, len(out)-1; i < j; i, j = i+1, j-1 {
		out[i], out[j] = out[j], out[i]
	}
	return out, nil
}

func (r *SampleRepository) FindByMonitorAndPeriod(ctx context.Context, monitorID string, from, to time.Time) ([]collector.Sample, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT "timestamp", success, latency_ns
		FROM samples
		WHERE monitor_id = $1 AND "timestamp" >= $2 AND "timestamp" < $3
		ORDER BY "timestamp"
	`, monitorID, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []collector.Sample
	for rows.Next() {
		s, err := scanSample(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, rows.Err()
}

func (r *SampleRepository) FindLastTimestamps(ctx context.Context, monitorIDs []string) (map[string]time.Time, error) {
	if len(monitorIDs) == 0 {
		return map[string]time.Time{}, nil
	}
	rows, err := r.pool.Query(ctx, `
		SELECT DISTINCT ON (monitor_id) monitor_id, "timestamp"
		FROM samples
		WHERE monitor_id = ANY($1)
		ORDER BY monitor_id, "timestamp" DESC
	`, monitorIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make(map[string]time.Time, len(monitorIDs))
	for rows.Next() {
		var monitorID string
		var ts time.Time
		if err := rows.Scan(&monitorID, &ts); err != nil {
			return nil, err
		}
		out[monitorID] = ts
	}
	return out, rows.Err()
}

func scanSample(row rowScanner) (collector.Sample, error) {
	var (
		s         collector.Sample
		latencyNS int64
	)

	if err := row.Scan(&s.Timestamp, &s.Success, &latencyNS); err != nil {
		return collector.Sample{}, err
	}

	s.Latency = time.Duration(latencyNS)
	return s, nil
}
