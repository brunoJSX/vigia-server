package stub

import (
	"context"
	"log"
)

type WhatsAppProvider struct {
	logger *log.Logger
}

func NewWhatsAppProvider(logger *log.Logger) *WhatsAppProvider {
	if logger == nil {
		logger = log.Default()
	}
	return &WhatsAppProvider{logger: logger}
}

func (p *WhatsAppProvider) Send(ctx context.Context, recipient, message string) error {
	p.logger.Printf("whatsapp [stub]: to=%s msg=%q", recipient, message)
	return nil
}
