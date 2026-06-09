package application

import (
	"context"
	"errors"
	"regexp"

	"github.com/vigia/vigia-v1/internal/account/account"
	"github.com/vigia/vigia-v1/internal/shared/clock"
)

var e164Re = regexp.MustCompile(`^\+[1-9]\d{6,14}$`)

var ErrInvalidPhoneNumber = errors.New("whatsapp_number must be in E.164 format (e.g. +5511999999999)")
var ErrAccountNotFound = errors.New("account not found")

type UpdateAccountInput struct {
	WhatsAppNumber string
}

type UpdateAccount struct {
	accounts AccountRepository
	clock    clock.Clock
}

func NewUpdateAccount(accounts AccountRepository, clk clock.Clock) *UpdateAccount {
	return &UpdateAccount{accounts: accounts, clock: clk}
}

func (uc *UpdateAccount) Execute(ctx context.Context, userID string, in UpdateAccountInput) (account.Account, error) {
	if in.WhatsAppNumber != "" && !e164Re.MatchString(in.WhatsAppNumber) {
		return account.Account{}, ErrInvalidPhoneNumber
	}

	a, err := uc.accounts.FindByID(ctx, userID)
	if err != nil {
		return account.Account{}, err
	}
	if a == nil {
		return account.Account{}, ErrAccountNotFound
	}

	a.WhatsAppNumber = in.WhatsAppNumber
	a.UpdatedAt = uc.clock()

	if err := uc.accounts.Save(ctx, *a); err != nil {
		return account.Account{}, err
	}
	return *a, nil
}
