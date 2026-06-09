package application

import (
	"context"
	"math"
	"time"

	"github.com/vigia/vigia-v1/internal/observability/incident"
	"github.com/vigia/vigia-v1/internal/shared/clock"
)

type QueryAggregateHistory struct {
	monitors  MonitorRepository
	incidents IncidentRepository
	clock     clock.Clock
}

func NewQueryAggregateHistory(monitors MonitorRepository, incidents IncidentRepository, clk clock.Clock) *QueryAggregateHistory {
	return &QueryAggregateHistory{monitors: monitors, incidents: incidents, clock: clk}
}

type DailyStats struct {
	Date            string  `json:"date"`
	IncidentsCount  int     `json:"incidents_count"`
	DowntimeSeconds float64 `json:"downtime_seconds"`
}

type AggregateHistoryResult struct {
	AvailabilityPercentage float64      `json:"availability_percentage"`
	TotalIncidents         int          `json:"total_incidents"`
	TotalDowntimeSeconds   float64      `json:"total_downtime_seconds"`
	Daily                  []DailyStats `json:"daily"`
}

func (uc *QueryAggregateHistory) Execute(ctx context.Context, days int) (AggregateHistoryResult, error) {
	now := uc.clock()
	to := now
	from := now.AddDate(0, 0, -days)

	monitors, err := uc.monitors.FindAll(ctx)
	if err != nil {
		return AggregateHistoryResult{}, err
	}

	incidents, err := uc.incidents.FindByPeriod(ctx, from, to)
	if err != nil {
		return AggregateHistoryResult{}, err
	}

	periodSecs := to.Sub(from).Seconds()

	// Per-monitor downtime for worst-case availability
	incidentsByMonitor := make(map[string][]incident.Incident, len(monitors))
	for _, inc := range incidents {
		incidentsByMonitor[inc.MonitorID] = append(incidentsByMonitor[inc.MonitorID], inc)
	}

	worstAvailability := 100.0
	for _, m := range monitors {
		downtime := totalDowntimeInInterval(incidentsByMonitor[m.ID], from, to, now)
		avail := (periodSecs - downtime) / periodSecs * 100
		if avail < worstAvailability {
			worstAvailability = avail
		}
	}

	// Total incidents opened within the period
	totalIncidents := 0
	for _, inc := range incidents {
		if !inc.OpenedAt.Before(from) {
			totalIncidents++
		}
	}

	// Total downtime across all monitors within the period
	totalDowntime := totalDowntimeInInterval(incidents, from, to, now)

	return AggregateHistoryResult{
		AvailabilityPercentage: math.Round(worstAvailability*100) / 100,
		TotalIncidents:         totalIncidents,
		TotalDowntimeSeconds:   math.Round(totalDowntime*100) / 100,
		Daily:                  buildDailyStats(incidents, from, to, now, days),
	}, nil
}

func totalDowntimeInInterval(incidents []incident.Incident, from, to, now time.Time) float64 {
	total := 0.0
	for _, inc := range incidents {
		total += incidentOverlapSeconds(inc, from, to, now)
	}
	return total
}

// incidentOverlapSeconds returns the seconds an incident overlapped with [from, to).
func incidentOverlapSeconds(inc incident.Incident, from, to, now time.Time) float64 {
	start := inc.OpenedAt
	if start.Before(from) {
		start = from
	}

	end := now
	if inc.ResolvedAt != nil {
		end = *inc.ResolvedAt
	}
	if end.After(to) {
		end = to
	}

	if !end.After(start) {
		return 0
	}
	return end.Sub(start).Seconds()
}

func buildDailyStats(incidents []incident.Incident, from, to, now time.Time, days int) []DailyStats {
	fromDay := time.Date(from.Year(), from.Month(), from.Day(), 0, 0, 0, 0, time.UTC)
	result := make([]DailyStats, 0, days)

	for d := 0; d < days; d++ {
		dayStart := fromDay.AddDate(0, 0, d)
		dayEnd := dayStart.Add(24 * time.Hour)

		// Clamp interval to period boundaries
		intStart := dayStart
		if intStart.Before(from) {
			intStart = from
		}
		intEnd := dayEnd
		if intEnd.After(to) {
			intEnd = to
		}

		count := 0
		downtime := 0.0
		for _, inc := range incidents {
			if !inc.OpenedAt.Before(dayStart) && inc.OpenedAt.Before(dayEnd) {
				count++
			}
			downtime += incidentOverlapSeconds(inc, intStart, intEnd, now)
		}

		result = append(result, DailyStats{
			Date:            dayStart.Format("2006-01-02"),
			IncidentsCount:  count,
			DowntimeSeconds: math.Round(downtime*100) / 100,
		})
	}

	return result
}
