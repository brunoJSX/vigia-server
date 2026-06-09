package notification_test

import (
	"testing"
	"time"

	"github.com/vigia/vigia-v1/internal/notification/notification"
)

func newNotification() notification.Notification {
	return notification.New("n-1", notification.TypeIncidentOpened, "+5511999999999",
		notification.Payload{MonitorName: "Checkout"}, time.Now())
}

// RN-N001: nasce com status Pending e Attempts zero.
func TestNotification_NewIsPending(t *testing.T) {
	n := newNotification()
	if n.Status != notification.StatusPending {
		t.Fatalf("expected Pending, got %q", n.Status)
	}
	if n.Attempts != 0 {
		t.Fatalf("expected Attempts 0, got %d", n.Attempts)
	}
}

// RN-N002: Delivered é estado final.
func TestNotification_MarkDelivered_IsFinal(t *testing.T) {
	n := newNotification()
	now := time.Now()
	n.MarkDelivered(now)

	if n.Status != notification.StatusDelivered {
		t.Fatalf("expected Delivered, got %q", n.Status)
	}
	if n.DeliveredAt == nil || !n.DeliveredAt.Equal(now) {
		t.Fatalf("expected DeliveredAt %v, got %v", now, n.DeliveredAt)
	}
	if n.IsActionable() {
		t.Fatal("Delivered notification must not be actionable (RN-N007)")
	}
}

// RN-N003 / RN-N004: falhas incrementam Attempts; ao atingir MaxAttempts → Dead.
func TestNotification_RecordFailure_TransitionsToDead(t *testing.T) {
	n := newNotification()

	for i := 1; i < notification.MaxAttempts; i++ {
		n.RecordFailure()
		if n.Status != notification.StatusFailed {
			t.Fatalf("attempt %d: expected Failed, got %q", i, n.Status)
		}
		if n.Attempts != i {
			t.Fatalf("attempt %d: expected Attempts %d, got %d", i, i, n.Attempts)
		}
		if !n.IsActionable() {
			t.Fatalf("attempt %d: should still be actionable", i)
		}
	}

	n.RecordFailure()
	if n.Status != notification.StatusDead {
		t.Fatalf("expected Dead after %d failures, got %q", notification.MaxAttempts, n.Status)
	}
	if n.IsActionable() {
		t.Fatal("Dead notification must not be actionable (RN-N007)")
	}
}

// Render: template correto por tipo.
func TestRender_IncidentOpened(t *testing.T) {
	n := notification.New("n-1", notification.TypeIncidentOpened, "+55",
		notification.Payload{MonitorName: "Checkout"}, time.Now())
	got := notification.Render(n)
	want := "🔴 Problema detectado — Checkout está fora do ar."
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestRender_IncidentResolved_WithDuration(t *testing.T) {
	n := notification.New("n-1", notification.TypeIncidentResolved, "+55",
		notification.Payload{MonitorName: "Checkout", Duration: "38m00s"}, time.Now())
	got := notification.Render(n)
	want := "✅ Problema resolvido — Checkout voltou ao normal. Duração: 38m00s."
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}
