package product

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Category string

const (
	CategoryTopup    Category = "topup"
	CategoryPulsa    Category = "pulsa"
	CategoryInternet Category = "internet"
)

type Product struct {
	bun.BaseModel `bun:"table:products"`

	ID        uuid.UUID `bun:",pk,type:uuid,default:gen_random_uuid()" json:"id"`
	Name      string    `json:"name"`
	Provider  string    `json:"provider"`
	Category  Category  `json:"category"`
	SKU       string    `json:"sku"`
	Price     int       `json:"price"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp" json:"updated_at"`
}

type Service struct {
	db *bun.DB
}

func NewService(db *bun.DB) *Service {
	return &Service{db: db}
}

func (s *Service) ListActive(ctx context.Context) ([]Product, error) {
	var products []Product
	err := s.db.NewSelect().
		Model(&products).
		Where("active = true").
		Order("category, price").
		Scan(ctx)
	return products, err
}

func (s *Service) FindByID(ctx context.Context, id uuid.UUID) (*Product, error) {
	var product Product
	err := s.db.NewSelect().Model(&product).Where("id = ?", id).Limit(1).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return &product, nil
}

// SeedDefaults inserts a curated set of SKUs for local/dev usage.
func (s *Service) SeedDefaults(ctx context.Context) error {
	initial := []Product{
		{Name: "Topup Gojek 50k", Provider: "Gojek", Category: CategoryTopup, SKU: "GJ-50", Price: 50000, Active: true},
		{Name: "Topup Grab 50k", Provider: "Grab", Category: CategoryTopup, SKU: "GR-50", Price: 50000, Active: true},
		{Name: "Pulsa Telkomsel 25k", Provider: "Telkomsel", Category: CategoryPulsa, SKU: "TSEL-25", Price: 25000, Active: true},
		{Name: "Paket Data 10GB", Provider: "XL", Category: CategoryInternet, SKU: "XL-10GB", Price: 45000, Active: true},
	}

	for _, p := range initial {
		p.ID = uuid.New()
		p.CreatedAt = time.Now()
		p.UpdatedAt = time.Now()
		_, _ = s.db.NewInsert().Model(&p).Ignore().Exec(ctx)
	}
	return nil
}
