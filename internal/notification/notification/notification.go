package notification

import (
	"fmt"
	"time"
)

type Status string

const (
	StatusPending   Status = "pending"
	StatusFailed    Status = "failed"
	StatusDelivered Status = "delivered"
	StatusDead      Status = "dead"
)

type Type string

const (
	TypeIncidentOpened   Type = "incident_opened"
	TypeIncidentResolved Type = "incident_resolved"
)

const MaxAttempts = 3

type Payload struct {
	MonitorName string `json:"monitor_name"`
	Duration    string `json:"duration,omitempty"`
}

type Notification struct {
	ID          string
	Type        Type
	Recipient   string
	Payload     Payload
	Status      Status
	Attempts    int
	CreatedAt   time.Time
	DeliveredAt *time.Time
}

func New(id string, t Type, recipient string, payload Payload, createdAt time.Time) Notification {
	return Notification{
		ID:        id,
		Type:      t,
		Recipient: recipient,
		Payload:   payload,
		Status:    StatusPending,
		Attempts:  0,
		CreatedAt: createdAt,
	}
}

// MarkDelivered transitions to Delivered — final state (RN-N002).
func (n *Notification) MarkDelivered(now time.Time) {
	n.Status = StatusDelivered
	n.DeliveredAt = &now
}

// RecordFailure increments Attempts and transitions to Failed or Dead (RN-N003, RN-N004).
func (n *Notification) RecordFailure() {
	n.Attempts++
	if n.Attempts >= MaxAttempts {
		n.Status = StatusDead
	} else {
		n.Status = StatusFailed
	}
}

// IsActionable reports whether the Worker should attempt delivery (RN-N007).
func (n Notification) IsActionable() bool {
	return (n.Status == StatusPending || n.Status == StatusFailed) && n.Attempts < MaxAttempts
}

// Render produces the WhatsApp message text from the notification type and payload.
func Render(n Notification) string {
	switch n.Type {
	case TypeIncidentOpened:
		return fmt.Sprintf("🔴 Problema detectado — %s está fora do ar.", n.Payload.MonitorName)
	case TypeIncidentResolved:
		if n.Payload.Duration != "" {
			return fmt.Sprintf("✅ Problema resolvido — %s voltou ao normal. Duração: %s.", n.Payload.MonitorName, n.Payload.Duration)
		}
		return fmt.Sprintf("✅ Problema resolvido — %s voltou ao normal.", n.Payload.MonitorName)
	}
	return ""
}
