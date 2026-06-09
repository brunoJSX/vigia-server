package memory

import (
	"context"
	"sync"

	"github.com/vigia/vigia-v1/internal/notification/notification"
)

type NotificationRepository struct {
	mu            sync.Mutex
	notifications map[string]notification.Notification
}

func NewNotificationRepository() *NotificationRepository {
	return &NotificationRepository{notifications: make(map[string]notification.Notification)}
}

func (r *NotificationRepository) Save(ctx context.Context, n notification.Notification) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.notifications[n.ID] = n
	return nil
}

func (r *NotificationRepository) FindActionable(ctx context.Context) ([]notification.Notification, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var out []notification.Notification
	for _, n := range r.notifications {
		if n.IsActionable() {
			out = append(out, n)
		}
	}
	return out, nil
}
