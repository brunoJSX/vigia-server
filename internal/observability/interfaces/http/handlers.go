// Package http exposes the Observability use cases over HTTP — stdlib
// net/http only, no framework (golang-conventions: prefer the simplest
// explicit code; there is no boundary here that would justify one).
//
// Surface is limited to what the client actually drives: Monitor management
// (the RN-037 gap the spec implies but doesn't name) and history queries.
// CheckMonitor, ResolveIncident and BuildDailySummary are pipeline/scheduled
// internals — not client-facing endpoints.
package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/vigia/vigia-v1/internal/observability/application"
	"github.com/vigia/vigia-v1/internal/observability/incident"
	"github.com/vigia/vigia-v1/internal/observability/monitor"
)

type Handlers struct {
	createMonitor         *application.CreateMonitor
	pauseMonitor          *application.PauseMonitor
	resumeMonitor         *application.ResumeMonitor
	disableMonitor        *application.DisableMonitor
	queryHistory          *application.QueryHistory
	queryMonitors         *application.QueryMonitors
	queryIncidents        *application.QueryIncidents
	queryAggregateHistory *application.QueryAggregateHistory
}

func NewHandlers(
	createMonitor *application.CreateMonitor,
	pauseMonitor *application.PauseMonitor,
	resumeMonitor *application.ResumeMonitor,
	disableMonitor *application.DisableMonitor,
	queryHistory *application.QueryHistory,
	queryMonitors *application.QueryMonitors,
	queryIncidents *application.QueryIncidents,
	queryAggregateHistory *application.QueryAggregateHistory,
) *Handlers {
	return &Handlers{
		createMonitor:         createMonitor,
		pauseMonitor:          pauseMonitor,
		resumeMonitor:         resumeMonitor,
		disableMonitor:        disableMonitor,
		queryHistory:          queryHistory,
		queryMonitors:         queryMonitors,
		queryIncidents:        queryIncidents,
		queryAggregateHistory: queryAggregateHistory,
	}
}

type createMonitorRequest struct {
	Name                          string  `json:"name"`
	Description                   string  `json:"description"`
	Target                        string  `json:"target"`
	Type                          string  `json:"type"`
	Threshold                     int     `json:"threshold"`
	IntervalSeconds               int     `json:"interval_seconds"`
	AcceptableResponseTimeSeconds float64 `json:"acceptable_response_time_seconds,omitempty"`
}

type monitorResponse struct {
	ID                            string  `json:"id"`
	Name                          string  `json:"name"`
	Description                   string  `json:"description"`
	Target                        string  `json:"target"`
	Type                          string  `json:"type"`
	Status                        string  `json:"status"`
	Threshold                     int     `json:"threshold"`
	IntervalSeconds               int     `json:"interval_seconds"`
	AcceptableResponseTimeSeconds float64 `json:"acceptable_response_time_seconds,omitempty"`
}

type monitorViewResponse struct {
	monitorResponse
	CurrentState  string  `json:"current_state"`
	LastCheckedAt *string `json:"last_checked_at"`
}

func newMonitorResponse(m monitor.Monitor) monitorResponse {
	return monitorResponse{
		ID:                            m.ID,
		Name:                          m.Name,
		Description:                   m.Description,
		Target:                        m.Target,
		Type:                          string(m.Type),
		Status:                        string(m.Status),
		Threshold:                     m.Threshold,
		IntervalSeconds:               int(m.Interval.Seconds()),
		AcceptableResponseTimeSeconds: m.AcceptableResponseTime.Seconds(),
	}
}

func newMonitorViewResponse(v application.MonitorView) monitorViewResponse {
	r := monitorViewResponse{
		monitorResponse: newMonitorResponse(v.Monitor),
		CurrentState:    string(v.CurrentState),
	}
	if v.LastCheckedAt != nil {
		s := v.LastCheckedAt.UTC().Format(time.RFC3339)
		r.LastCheckedAt = &s
	}
	return r
}

func (h *Handlers) CreateMonitor(w http.ResponseWriter, r *http.Request) {
	var req createMonitorRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}
	if monitor.Type(req.Type) == monitor.TypeCheckout && req.AcceptableResponseTimeSeconds <= 0 {
		http.Error(w, "acceptable_response_time_seconds is required for checkout monitors", http.StatusBadRequest)
		return
	}

	m, err := h.createMonitor.Execute(r.Context(), application.CreateMonitorInput{
		Name:                   req.Name,
		Description:            req.Description,
		Target:                 req.Target,
		Type:                   monitor.Type(req.Type),
		Threshold:              req.Threshold,
		Interval:               time.Duration(req.IntervalSeconds) * time.Second,
		AcceptableResponseTime: time.Duration(req.AcceptableResponseTimeSeconds * float64(time.Second)),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, newMonitorResponse(m))
}

func (h *Handlers) ListMonitors(w http.ResponseWriter, r *http.Request) {
	views, err := h.queryMonitors.Execute(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := make([]monitorViewResponse, 0, len(views))
	for _, v := range views {
		resp = append(resp, newMonitorViewResponse(v))
	}
	writeJSON(w, http.StatusOK, map[string]any{"monitors": resp})
}

func (h *Handlers) ListIncidents(w http.ResponseWriter, r *http.Request) {
	statusParam := r.URL.Query().Get("status")

	var status incident.Status
	switch statusParam {
	case "open":
		status = incident.StatusOpen
	case "resolved":
		status = incident.StatusResolved
	case "":
		status = ""
	default:
		http.Error(w, "status must be 'open', 'resolved', or omitted", http.StatusBadRequest)
		return
	}

	views, err := h.queryIncidents.Execute(r.Context(), application.QueryIncidentsInput{
		Status: status,
		Limit:  50,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	type incidentViewResponse struct {
		ID              string     `json:"id"`
		DisplayID       string     `json:"display_id"`
		MonitorID       string     `json:"monitor_id"`
		MonitorName     string     `json:"monitor_name"`
		Status          string     `json:"status"`
		OpenedAt        time.Time  `json:"opened_at"`
		ResolvedAt      *time.Time `json:"resolved_at,omitempty"`
		DurationSeconds float64    `json:"duration_seconds"`
	}

	resp := make([]incidentViewResponse, 0, len(views))
	for _, v := range views {
		resp = append(resp, incidentViewResponse{
			ID:              v.Incident.ID,
			DisplayID:       fmt.Sprintf("INC-%d", v.Incident.SequenceNumber),
			MonitorID:       v.Incident.MonitorID,
			MonitorName:     v.MonitorName,
			Status:          string(v.Incident.Status),
			OpenedAt:        v.Incident.OpenedAt,
			ResolvedAt:      v.Incident.ResolvedAt,
			DurationSeconds: v.Incident.Duration().Seconds(),
		})
	}
	writeJSON(w, http.StatusOK, map[string]any{"incidents": resp})
}

func (h *Handlers) PauseMonitor(w http.ResponseWriter, r *http.Request) {
	if err := h.pauseMonitor.Execute(r.Context(), r.PathValue("id")); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handlers) ResumeMonitor(w http.ResponseWriter, r *http.Request) {
	if err := h.resumeMonitor.Execute(r.Context(), r.PathValue("id")); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handlers) DisableMonitor(w http.ResponseWriter, r *http.Request) {
	if err := h.disableMonitor.Execute(r.Context(), r.PathValue("id")); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

type incidentResponse struct {
	ID              string     `json:"id"`
	MonitorID       string     `json:"monitor_id"`
	Status          string     `json:"status"`
	OpenedAt        time.Time  `json:"opened_at"`
	ResolvedAt      *time.Time `json:"resolved_at,omitempty"`
	DurationSeconds float64    `json:"duration_seconds"`
}

type historyResponse struct {
	Incidents              []incidentResponse `json:"incidents"`
	AvailabilityPercentage float64            `json:"availability_percentage"`
}

func newHistoryResponse(r application.QueryHistoryResult) historyResponse {
	incidents := make([]incidentResponse, 0, len(r.Incidents))
	for _, i := range r.Incidents {
		incidents = append(incidents, incidentResponse{
			ID:              i.ID,
			MonitorID:       i.MonitorID,
			Status:          string(i.Status),
			OpenedAt:        i.OpenedAt,
			ResolvedAt:      i.ResolvedAt,
			DurationSeconds: i.Duration().Seconds(),
		})
	}
	return historyResponse{Incidents: incidents, AvailabilityPercentage: r.AvailabilityPercentage}
}

func (h *Handlers) MonitorHistory(w http.ResponseWriter, r *http.Request) {
	from, err := time.Parse(time.RFC3339, r.URL.Query().Get("from"))
	if err != nil {
		http.Error(w, "invalid or missing 'from' (expected RFC3339)", http.StatusBadRequest)
		return
	}
	to, err := time.Parse(time.RFC3339, r.URL.Query().Get("to"))
	if err != nil {
		http.Error(w, "invalid or missing 'to' (expected RFC3339)", http.StatusBadRequest)
		return
	}

	result, err := h.queryHistory.Execute(r.Context(), application.QueryHistoryInput{
		MonitorID: r.PathValue("id"),
		From:      from,
		To:        to,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, newHistoryResponse(result))
}

func (h *Handlers) AggregateHistory(w http.ResponseWriter, r *http.Request) {
	days := 30
	if p := r.URL.Query().Get("days"); p != "" {
		if n, err := parseInt(p); err == nil && n > 0 && n <= 365 {
			days = n
		} else {
			http.Error(w, "days must be a positive integer up to 365", http.StatusBadRequest)
			return
		}
	}

	result, err := h.queryAggregateHistory.Execute(r.Context(), days)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, result)
}

func parseInt(s string) (int, error) {
	n := 0
	for _, c := range s {
		if c < '0' || c > '9' {
			return 0, fmt.Errorf("not a number")
		}
		n = n*10 + int(c-'0')
	}
	return n, nil
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
