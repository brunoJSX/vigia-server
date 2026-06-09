package application

import "context"

// PauseMonitor transitions a Monitor to Paused (RN-037).
type PauseMonitor struct {
	monitors MonitorRepository
}

func NewPauseMonitor(monitors MonitorRepository) *PauseMonitor {
	return &PauseMonitor{monitors: monitors}
}

func (uc *PauseMonitor) Execute(ctx context.Context, monitorID string) error {
	m, err := uc.monitors.FindByID(ctx, monitorID)
	if err != nil {
		return err
	}

	m.Pause()

	return uc.monitors.Save(ctx, m)
}
