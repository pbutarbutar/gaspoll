package mitra

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Mitra struct {
	bun.BaseModel `bun:"table:mitras"`

	ID        uuid.UUID `bun:",pk,type:uuid,default:gen_random_uuid()" json:"id"`
	Name      string    `json:"name"`
	Address   string    `json:"address"`
	Phone     string    `json:"phone"`
	CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp" json:"created_at"`
}

type Promo struct {
	bun.BaseModel `bun:"table:promos"`

	ID          uuid.UUID `bun:",pk,type:uuid,default:gen_random_uuid()" json:"id"`
	MitraID     uuid.UUID `json:"mitra_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Price       int       `json:"price"`
	CreatedAt   time.Time `bun:",nullzero,notnull,default:current_timestamp" json:"created_at"`
}

type Service struct {
	db *bun.DB
}

func NewService(db *bun.DB) *Service {
	return &Service{db: db}
}

func (s *Service) ListPromos(ctx context.Context, mitraID uuid.UUID) ([]Promo, error) {
	var promos []Promo
	err := s.db.NewSelect().
		Model(&promos).
		Where("mitra_id = ?", mitraID).
		Order("created_at DESC").
		Scan(ctx)
	return promos, err
}
