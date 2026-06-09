package application

import "context"

type ResolveRecipient struct {
	accounts AccountRepository
}

func NewResolveRecipient(accounts AccountRepository) *ResolveRecipient {
	return &ResolveRecipient{accounts: accounts}
}

func (uc *ResolveRecipient) Execute(ctx context.Context, accountID string) (string, error) {
	return uc.accounts.FindWhatsAppNumber(ctx, accountID)
}
