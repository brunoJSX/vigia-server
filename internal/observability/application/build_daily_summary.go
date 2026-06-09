package application

import (
	"context"
	"time"

	"github.com/vigia/vigia-v1/internal/observability/incident"
	"github.com/vigia/vigia-v1/internal/shared/clock"
)

// BuildDailySummary is the workflow BuildDailySummary — a daily scheduled
// run that consolidates the previous day's incidents and collected data per
// Monitor and makes a summary available to the client.
//
// DailySummary is not a domain concept (PA-? open question on its
// ownership/shape) — it is a result DTO this use case produces and emits.
type BuildDailySummary struct {
	monitors      MonitorRepository
	incidents     IncidentRepository
	samples       SampleRepository
	notifications NotificationPublisher
	clock         clock.Clock
}

func NewBuildDailySummary(
	monitors MonitorRepository,
	incidents IncidentRepository,
	samples SampleRepository,
	notifications NotificationPublisher,
	c clock.Clock,
) *BuildDailySummary {
	return &BuildDailySummary{
		monitors:      monitors,
		incidents:     incidents,
		samples:       samples,
		notifications: notifications,
		clock:         c,
	}
}

type DailySummaryEntry struct {
	MonitorID              string
	Incidents              []incident.Incident
	AvailabilityPercentage float64
}

type DailySummary struct {
	Date    time.Time
	Entries []DailySummaryEntry
}

func (uc *BuildDailySummary) Execute(ctx context.Context) (DailySummary, error) {
	now := uc.clock()
	dayStart := now.AddDate(0, 0, -1).Truncate(24 * time.Hour)
	dayEnd := dayStart.Add(24 * time.Hour)

	monitors, err := uc.monitors.FindActive(ctx)
	if err != nil {
		return DailySummary{}, err
	}

	summary := DailySummary{Date: dayStart}

	for _, m := range monitors {
		incidents, err := uc.incidents.FindByMonitorAndPeriod(ctx, m.ID, dayStart, dayEnd)
		if err != nil {
			return DailySummary{}, err
		}

		samples, err := uc.samples.FindByMonitorAndPeriod(ctx, m.ID, dayStart, dayEnd)
		if err != nil {
			return DailySummary{}, err
		}

		summary.Entries = append(summary.Entries, DailySummaryEntry{
			MonitorID:              m.ID,
			Incidents:              incidents,
			AvailabilityPercentage: availability(samples),
		})
	}

	if err := uc.notifications.Publish(ctx, Event{
		Kind:       EventDailySummaryReady,
		OccurredAt: now,
		Payload:    summary,
	}); err != nil {
		return DailySummary{}, err
	}

	return summary, nil
}
