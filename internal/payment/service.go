package payment

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Channel string

const (
	ChannelQRIS    Channel = "qris"
	ChannelVA      Channel = "va"
	ChannelEWallet Channel = "ewallet"
)

type Status string

const (
	StatusPending Status = "pending"
	StatusPaid    Status = "paid"
	StatusFailed  Status = "failed"
)

type Invoice struct {
	OrderID      uuid.UUID
	Amount       int
	Channel      Channel
	CheckoutURL  string
	ReferenceID  string
	ExpiredAt    time.Time
	CustomerNote string
}

// Service fakes a PSP/backoffice connection. Replace with real connectors later.
type Service struct {
	baseURL string
}

func NewService(baseURL string) *Service {
	return &Service{baseURL: baseURL}
}

func (s *Service) CreateInvoice(ctx context.Context, orderID uuid.UUID, amount int, channel Channel) (*Invoice, error) {
	ref := orderID.String()
	return &Invoice{
		OrderID:     orderID,
		Amount:      amount,
		Channel:     channel,
		CheckoutURL: fmt.Sprintf("%s/pay/%s", s.baseURL, ref),
		ReferenceID: ref,
		ExpiredAt:   time.Now().Add(15 * time.Minute),
	}, nil
}

// VerifyWebhook should validate incoming webhook payloads; stubbed for now.
func (s *Service) VerifyWebhook(_ context.Context, ref string) (uuid.UUID, Status, error) {
	// In integration, verify signature then map ref -> orderID.
	orderID, err := uuid.Parse(ref)
	if err != nil {
		return uuid.Nil, StatusFailed, err
	}
	return orderID, StatusPaid, nil
}
