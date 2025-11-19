package settlement

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Settlement struct {
	bun.BaseModel `bun:"table:settlements"`

	ID         uuid.UUID `bun:",pk,type:uuid,default:gen_random_uuid()" json:"id"`
	MitraID    uuid.UUID `json:"mitra_id"`
	Amount     int       `json:"amount"`
	PeriodFrom time.Time `json:"period_from"`
	PeriodTo   time.Time `json:"period_to"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `bun:",nullzero,notnull,default:current_timestamp" json:"created_at"`
}

type Service struct {
	db *bun.DB
}

func NewService(db *bun.DB) *Service {
	return &Service{db: db}
}

func (s *Service) List(ctx context.Context) ([]Settlement, error) {
	var settlements []Settlement
	err := s.db.NewSelect().
		Model(&settlements).
		Order("period_from DESC").
		Scan(ctx)
	return settlements, err
}
