package application

import (
	"context"
	"fmt"
)

// PauseMonitor transitions a Monitor to Paused (RN-037).
type PauseMonitor struct {
	monitors MonitorRepository
}

func NewPauseMonitor(monitors MonitorRepository) *PauseMonitor {
	return &PauseMonitor{monitors: monitors}
}

func (uc *PauseMonitor) Execute(ctx context.Context, monitorID, accountID string) error {
	m, err := uc.monitors.FindByID(ctx, monitorID)
	if err != nil {
		return err
	}
	if m.AccountID != accountID {
		return fmt.Errorf("monitor %q not found", monitorID)
	}

	m.Pause()

	return uc.monitors.Save(ctx, m)
}
