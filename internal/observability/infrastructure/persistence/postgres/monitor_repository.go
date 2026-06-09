// Package postgres provides the Postgres-backed MonitorRepository /
// IncidentRepository / SampleRepository — explicit SQL via pgx, no ORM
// (golang-conventions: prefer explicit code over abstractions). Schema in
// schema.sql.
package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/vigia/vigia-v1/internal/observability/monitor"
)

type MonitorRepository struct {
	pool *pgxpool.Pool
}

func NewMonitorRepository(pool *pgxpool.Pool) *MonitorRepository {
	return &MonitorRepository{pool: pool}
}

const monitorColumns = `id, name, description, target, type, status, threshold, interval_ns, acceptable_response_time_ns`

func (r *MonitorRepository) Save(ctx context.Context, m monitor.Monitor) error {
	var art *int64
	if m.AcceptableResponseTime > 0 {
		v := int64(m.AcceptableResponseTime)
		art = &v
	}
	_, err := r.pool.Exec(ctx, `
		INSERT INTO monitors (id, name, description, target, type, status, threshold, interval_ns, acceptable_response_time_ns, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, now(), now())
		ON CONFLICT (id) DO UPDATE SET
			name                        = EXCLUDED.name,
			description                 = EXCLUDED.description,
			target                      = EXCLUDED.target,
			type                        = EXCLUDED.type,
			status                      = EXCLUDED.status,
			threshold                   = EXCLUDED.threshold,
			interval_ns                 = EXCLUDED.interval_ns,
			acceptable_response_time_ns = EXCLUDED.acceptable_response_time_ns,
			updated_at                  = now()
	`, m.ID, m.Name, m.Description, m.Target, string(m.Type), string(m.Status), m.Threshold, int64(m.Interval), art)
	return err
}

func (r *MonitorRepository) FindByID(ctx context.Context, id string) (monitor.Monitor, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT `+monitorColumns+` FROM monitors WHERE id = $1
	`, id)

	m, err := scanMonitor(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return monitor.Monitor{}, fmt.Errorf("monitor %q not found", id)
	}
	return m, err
}

func (r *MonitorRepository) FindActive(ctx context.Context) ([]monitor.Monitor, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT `+monitorColumns+` FROM monitors WHERE status = $1
	`, string(monitor.StatusActive))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []monitor.Monitor
	for rows.Next() {
		m, err := scanMonitor(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, m)
	}
	return out, rows.Err()
}

func (r *MonitorRepository) FindAll(ctx context.Context) ([]monitor.Monitor, error) {
	rows, err := r.pool.Query(ctx, `SELECT `+monitorColumns+` FROM monitors ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []monitor.Monitor
	for rows.Next() {
		m, err := scanMonitor(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, m)
	}
	return out, rows.Err()
}

type rowScanner interface {
	Scan(dest ...any) error
}

func scanMonitor(row rowScanner) (monitor.Monitor, error) {
	var (
		m           monitor.Monitor
		monitorType string
		status      string
		intervalNS  int64
		artNS       *int64
	)

	if err := row.Scan(&m.ID, &m.Name, &m.Description, &m.Target, &monitorType, &status, &m.Threshold, &intervalNS, &artNS); err != nil {
		return monitor.Monitor{}, err
	}

	m.Type = monitor.Type(monitorType)
	m.Status = monitor.Status(status)
	m.Interval = time.Duration(intervalNS)
	if artNS != nil {
		m.AcceptableResponseTime = time.Duration(*artNS)
	}
	return m, nil
}
