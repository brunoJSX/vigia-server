package application

import (
	"context"
	"fmt"
)

// ResumeMonitor transitions a Monitor back to Active (RN-037).
type ResumeMonitor struct {
	monitors MonitorRepository
}

func NewResumeMonitor(monitors MonitorRepository) *ResumeMonitor {
	return &ResumeMonitor{monitors: monitors}
}

func (uc *ResumeMonitor) Execute(ctx context.Context, monitorID, accountID string) error {
	m, err := uc.monitors.FindByID(ctx, monitorID)
	if err != nil {
		return err
	}
	if m.AccountID != accountID {
		return fmt.Errorf("monitor %q not found", monitorID)
	}

	m.Resume()

	return uc.monitors.Save(ctx, m)
}
