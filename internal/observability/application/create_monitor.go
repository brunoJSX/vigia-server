package application

import (
	"context"
	"time"

	"github.com/vigia/vigia-v1/internal/observability/monitor"
	"github.com/vigia/vigia-v1/internal/shared/id"
)

// CreateMonitor registers a new Monitor configuration. monitor.New always
// starts it Active (RN-037) — there is no "draft" or "uninitialized" status.
type CreateMonitor struct {
	monitors MonitorRepository
	ids      id.Generator
}

func NewCreateMonitor(monitors MonitorRepository, ids id.Generator) *CreateMonitor {
	return &CreateMonitor{monitors: monitors, ids: ids}
}

type CreateMonitorInput struct {
	AccountID              string
	Name                   string
	Description            string
	Target                 string
	Type                   monitor.Type
	Threshold              int
	Interval               time.Duration
	AcceptableResponseTime time.Duration
}

func (uc *CreateMonitor) Execute(ctx context.Context, in CreateMonitorInput) (monitor.Monitor, error) {
	m := monitor.New(uc.ids(), in.AccountID, in.Name, in.Description, in.Target, in.Type, in.Threshold, in.Interval, in.AcceptableResponseTime)

	if err := uc.monitors.Save(ctx, m); err != nil {
		return monitor.Monitor{}, err
	}

	return m, nil
}
