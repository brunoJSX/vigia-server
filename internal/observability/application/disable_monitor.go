package application

import "context"

// DisableMonitor transitions a Monitor to Disabled (RN-037) — the
// configuration stops being checked, but is kept for history.
type DisableMonitor struct {
	monitors MonitorRepository
}

func NewDisableMonitor(monitors MonitorRepository) *DisableMonitor {
	return &DisableMonitor{monitors: monitors}
}

func (uc *DisableMonitor) Execute(ctx context.Context, monitorID string) error {
	m, err := uc.monitors.FindByID(ctx, monitorID)
	if err != nil {
		return err
	}

	m.Disable()

	return uc.monitors.Save(ctx, m)
}
