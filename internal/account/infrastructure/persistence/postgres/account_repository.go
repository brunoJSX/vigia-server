package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/vigia/vigia-v1/internal/account/account"
)

type AccountRepository struct {
	pool *pgxpool.Pool
}

func NewAccountRepository(pool *pgxpool.Pool) *AccountRepository {
	return &AccountRepository{pool: pool}
}

func (r *AccountRepository) FindByID(ctx context.Context, id string) (*account.Account, error) {
	var (
		a              account.Account
		whatsappNumber *string
	)
	err := r.pool.QueryRow(ctx, `
		SELECT id, whatsapp_number, created_at, updated_at
		FROM accounts WHERE id = $1
	`, id).Scan(&a.ID, &whatsappNumber, &a.CreatedAt, &a.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if whatsappNumber != nil {
		a.WhatsAppNumber = *whatsappNumber
	}
	return &a, nil
}

func (r *AccountRepository) Save(ctx context.Context, a account.Account) error {
	var whatsappNumber *string
	if a.WhatsAppNumber != "" {
		whatsappNumber = &a.WhatsAppNumber
	}
	_, err := r.pool.Exec(ctx, `
		INSERT INTO accounts (id, whatsapp_number, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (id) DO UPDATE SET
			whatsapp_number = EXCLUDED.whatsapp_number,
			updated_at      = EXCLUDED.updated_at
	`, a.ID, whatsappNumber, a.CreatedAt, a.UpdatedAt)
	return err
}

func (r *AccountRepository) FindWhatsAppNumber(ctx context.Context) (string, error) {
	var number string
	err := r.pool.QueryRow(ctx, `
		SELECT whatsapp_number FROM accounts
		WHERE whatsapp_number IS NOT NULL
		LIMIT 1
	`).Scan(&number)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", nil
	}
	return number, err
}

