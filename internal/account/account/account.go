package account

import "time"

type Account struct {
	ID             string
	WhatsAppNumber string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func New(id string, now time.Time) Account {
	return Account{ID: id, CreatedAt: now, UpdatedAt: now}
}
