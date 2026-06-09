package monitor

import "time"

// Status represents the operational status of a Monitor (RN-037).
type Status string

const (
	StatusActive   Status = "active"
	StatusPaused   Status = "paused"
	StatusDisabled Status = "disabled"
)

// Type represents what kind of operation a Monitor watches.
// Checkout is a specialized Monitor type (RN-024); modeled as a plain value
// until PA-001 (whether Uptime/Checkout/Dependency deserve their own structs)
// is resolved.
type Type string

const (
	TypeUptime     Type = "uptime"
	TypeCheckout   Type = "checkout"
	TypeDependency Type = "dependency"
)

// Monitor is configuration: it defines what should be observed, how often,
// and the thresholds used to interpret results. It does not run checks
// itself — the processing pipeline (Collector/Analyzer) consumes this
// configuration to do that.
//
// AcceptableResponseTime is a Checkout-specific field (RN-025): the latency
// threshold above which a response counts as Lentidão. Zero for Uptime and
// Dependency — those types have no domain-level time threshold.
type Monitor struct {
	ID                     string
	AccountID              string
	Name                   string
	Description            string
	Target                 string
	Type                   Type
	Status                 Status
	Threshold              int
	Interval               time.Duration
	AcceptableResponseTime time.Duration
}

// New creates a Monitor with status Active — a Monitor always has a status
// (RN-037), so there is no constructor path that leaves it undefined.
// For Checkout monitors, acceptableResponseTime must be > 0 (RN-025).
// For other types, pass 0.
func New(id, accountID, name, description, target string, monitorType Type, threshold int, interval, acceptableResponseTime time.Duration) Monitor {
	return Monitor{
		ID:                     id,
		AccountID:              accountID,
		Name:                   name,
		Description:            description,
		Target:                 target,
		Type:                   monitorType,
		Status:                 StatusActive,
		Threshold:              threshold,
		Interval:               interval,
		AcceptableResponseTime: acceptableResponseTime,
	}
}

func (m *Monitor) Pause() {
	m.Status = StatusPaused
}

func (m *Monitor) Resume() {
	m.Status = StatusActive
}

func (m *Monitor) Disable() {
	m.Status = StatusDisabled
}

func (m Monitor) IsActive() bool {
	return m.Status == StatusActive
}
