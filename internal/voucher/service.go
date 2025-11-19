package voucher

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/skip2/go-qrcode"
	"github.com/uptrace/bun"
)

type Status string

const (
	StatusIssued   Status = "issued"
	StatusRedeemed Status = "redeemed"
	StatusExpired  Status = "expired"
)

type Voucher struct {
	bun.BaseModel `bun:"table:vouchers"`

	ID         uuid.UUID `bun:",pk,type:uuid,default:gen_random_uuid()" json:"id"`
	UserID     uuid.UUID `json:"user_id"`
	Code       string    `json:"code"`
	Value      int       `json:"value"`
	Partner    string    `json:"partner"`
	Status     Status    `json:"status"`
	ExpiresAt  time.Time `json:"expires_at"`
	RedeemedAt time.Time `json:"redeemed_at"`
	CreatedAt  time.Time `bun:",nullzero,notnull,default:current_timestamp" json:"created_at"`
}

type Service struct {
	db *bun.DB
}

func NewService(db *bun.DB) *Service {
	return &Service{db: db}
}

func (s *Service) Issue(ctx context.Context, userID uuid.UUID, value int, partner string) (*Voucher, error) {
	code := fmt.Sprintf("GSP-%s", uuid.New().String())
	voucher := &Voucher{
		ID:        uuid.New(),
		UserID:    userID,
		Code:      code,
		Value:     value,
		Partner:   partner,
		Status:    StatusIssued,
		ExpiresAt: time.Now().Add(45 * 24 * time.Hour),
		CreatedAt: time.Now(),
	}

	_, err := s.db.NewInsert().Model(voucher).Exec(ctx)
	if err != nil {
		return nil, err
	}
	return voucher, nil
}

func (s *Service) Redeem(ctx context.Context, code string, partner string) error {
	_, err := s.db.NewUpdate().
		Model((*Voucher)(nil)).
		Set("status = ?", StatusRedeemed).
		Set("partner = ?", partner).
		Set("redeemed_at = ?", time.Now()).
		Where("code = ? AND status = ?", code, StatusIssued).
		Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) ListByUser(ctx context.Context, userID uuid.UUID) ([]Voucher, error) {
	var vouchers []Voucher
	err := s.db.NewSelect().
		Model(&vouchers).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Scan(ctx)
	return vouchers, err
}

func (v *Voucher) QRDataURI() string {
	png, err := qrcode.Encode(v.Code, qrcode.Medium, 256)
	if err != nil {
		log.Error().Err(err).Msg("unable to generate QR code")
		return ""
	}
	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(png)
}
