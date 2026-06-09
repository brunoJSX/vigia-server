package application

import (
	"context"

	"github.com/vigia/vigia-v1/internal/account/account"
	"github.com/vigia/vigia-v1/internal/shared/clock"
)

type GetAccount struct {
	accounts AccountRepository
	clock    clock.Clock
}

func NewGetAccount(accounts AccountRepository, clk clock.Clock) *GetAccount {
	return &GetAccount{accounts: accounts, clock: clk}
}

func (uc *GetAccount) Execute(ctx context.Context, userID string) (account.Account, error) {
	a, err := uc.accounts.FindByID(ctx, userID)
	if err != nil {
		return account.Account{}, err
	}
	if a != nil {
		return *a, nil
	}
	created := account.New(userID, uc.clock())
	if err := uc.accounts.Save(ctx, created); err != nil {
		return account.Account{}, err
	}
	return created, nil
}
