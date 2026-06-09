package application

import (
	"context"

	"github.com/vigia/vigia-v1/internal/account/account"
)

type AccountRepository interface {
	FindByID(ctx context.Context, id string) (*account.Account, error)
	Save(ctx context.Context, a account.Account) error
	FindWhatsAppNumber(ctx context.Context, accountID string) (string, error)
}
