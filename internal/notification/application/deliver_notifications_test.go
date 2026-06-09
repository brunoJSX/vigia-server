package application_test

import (
	"context"
	"errors"
	"testing"
	"time"

	notifapp "github.com/vigia/vigia-v1/internal/notification/application"
	notifmemory "github.com/vigia/vigia-v1/internal/notification/infrastructure/persistence/memory"
	"github.com/vigia/vigia-v1/internal/notification/notification"
)

type fakeProvider struct {
	err   error
	calls int
}

func (p *fakeProvider) Send(_ context.Context, _, _ string) error {
	p.calls++
	return p.err
}

func seedNotification(t *testing.T, repo *notifmemory.NotificationRepository, n notification.Notification) {
	t.Helper()
	if err := repo.Save(context.Background(), n); err != nil {
		t.Fatalf("seed: %v", err)
	}
}

func newPending(id string) notification.Notification {
	return notification.New(id, notification.TypeIncidentOpened, "+5511999999999",
		notification.Payload{MonitorName: "Checkout"}, time.Now())
}

// DeliverNotifications: sucesso → status Delivered.
func TestDeliverNotifications_Success(t *testing.T) {
	ctx := context.Background()
	repo := notifmemory.NewNotificationRepository()
	provider := &fakeProvider{}
	uc := notifapp.NewDeliverNotifications(repo, provider, time.Now)

	seedNotification(t, repo, newPending("n-1"))

	if err := uc.Execute(ctx); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	pending, _ := repo.FindActionable(ctx)
	if len(pending) != 0 {
		t.Fatalf("expected no actionable notifications after delivery, got %d", len(pending))
	}
	if provider.calls != 1 {
		t.Fatalf("expected 1 provider call, got %d", provider.calls)
	}
}

// RN-N004: falhas acumuladas até MaxAttempts → Dead.
func TestDeliverNotifications_RetriesUntilDead(t *testing.T) {
	ctx := context.Background()
	repo := notifmemory.NewNotificationRepository()
	provider := &fakeProvider{err: errors.New("provider down")}
	uc := notifapp.NewDeliverNotifications(repo, provider, time.Now)

	seedNotification(t, repo, newPending("n-1"))

	for i := 0; i < notification.MaxAttempts; i++ {
		if err := uc.Execute(ctx); err != nil {
			t.Fatalf("unexpected error on attempt %d: %v", i+1, err)
		}
	}

	pending, _ := repo.FindActionable(ctx)
	if len(pending) != 0 {
		t.Fatalf("expected no actionable notifications after exhausting retries, got %d", len(pending))
	}
	if provider.calls != notification.MaxAttempts {
		t.Fatalf("expected %d provider calls, got %d", notification.MaxAttempts, provider.calls)
	}
}

// RN-N006: sem destinatário → nenhuma Notification criada.
func TestEnqueueNotification_EmptyRecipient_DoesNothing(t *testing.T) {
	ctx := context.Background()
	repo := notifmemory.NewNotificationRepository()
	uc := notifapp.NewEnqueueNotification(repo, func() string { return "n-1" }, time.Now)

	err := uc.Execute(ctx, notifapp.EnqueueInput{
		Type:      notification.TypeIncidentOpened,
		Recipient: "",
		Payload:   notification.Payload{MonitorName: "Checkout"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	pending, _ := repo.FindActionable(ctx)
	if len(pending) != 0 {
		t.Fatalf("expected no notifications for empty recipient (RN-N006), got %d", len(pending))
	}
}
