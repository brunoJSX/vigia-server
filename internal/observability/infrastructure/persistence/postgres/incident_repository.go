package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/vigia/vigia-v1/internal/observability/incident"
)

type IncidentRepository struct {
	pool *pgxpool.Pool
}

func NewIncidentRepository(pool *pgxpool.Pool) *IncidentRepository {
	return &IncidentRepository{pool: pool}
}

func (r *IncidentRepository) Save(ctx context.Context, i incident.Incident) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO incidents (id, monitor_id, status, opened_at, resolved_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (id) DO UPDATE SET
			status      = EXCLUDED.status,
			resolved_at = EXCLUDED.resolved_at
	`, i.ID, i.MonitorID, string(i.Status), i.OpenedAt, i.ResolvedAt)
	return err
}

// FindOpenByMonitorID returns nil when there is none — RN-002 guarantees
// at most one Open Incident per Monitor.
func (r *IncidentRepository) FindOpenByMonitorID(ctx context.Context, monitorID string) (*incident.Incident, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, monitor_id, status, opened_at, resolved_at
		FROM incidents
		WHERE monitor_id = $1 AND status = $2
		LIMIT 1
	`, monitorID, string(incident.StatusOpen))

	i, err := scanIncident(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &i, nil
}

func (r *IncidentRepository) FindByMonitorAndPeriod(ctx context.Context, monitorID string, from, to time.Time) ([]incident.Incident, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, monitor_id, status, opened_at, resolved_at
		FROM incidents
		WHERE monitor_id = $1 AND opened_at >= $2 AND opened_at < $3
		ORDER BY opened_at
	`, monitorID, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []incident.Incident
	for rows.Next() {
		i, err := scanIncident(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, i)
	}
	return out, rows.Err()
}

func scanIncident(row rowScanner) (incident.Incident, error) {
	var (
		i          incident.Incident
		status     string
		resolvedAt *time.Time
	)

	if err := row.Scan(&i.ID, &i.MonitorID, &status, &i.OpenedAt, &resolvedAt); err != nil {
		return incident.Incident{}, err
	}

	i.Status = incident.Status(status)
	i.ResolvedAt = resolvedAt
	return i, nil
}
