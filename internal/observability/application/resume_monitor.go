package application

import "context"

// ResumeMonitor transitions a Monitor back to Active (RN-037).
type ResumeMonitor struct {
	monitors MonitorRepository
}

func NewResumeMonitor(monitors MonitorRepository) *ResumeMonitor {
	return &ResumeMonitor{monitors: monitors}
}

func (uc *ResumeMonitor) Execute(ctx context.Context, monitorID string) error {
	m, err := uc.monitors.FindByID(ctx, monitorID)
	if err != nil {
		return err
	}

	m.Resume()

	return uc.monitors.Save(ctx, m)
}
