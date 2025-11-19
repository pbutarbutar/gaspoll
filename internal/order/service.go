package order

import (
	"context"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/uptrace/bun"

	"gaspoll/internal/payment"
	"gaspoll/internal/product"
	"gaspoll/internal/reward"
	"gaspoll/internal/user"
)

type Status string

const (
	StatusPending    Status = "pending"
	StatusPaid       Status = "paid"
	StatusProcessing Status = "processing"
	StatusCompleted  Status = "completed"
	StatusFailed     Status = "failed"
)

type Order struct {
	bun.BaseModel `bun:"table:orders"`

	ID         uuid.UUID      `bun:",pk,type:uuid,default:gen_random_uuid()" json:"id"`
	UserID     uuid.UUID      `json:"user_id"`
	ProductID  uuid.UUID      `json:"product_id"`
	Channel    payment.Channel`json:"channel"`
	Amount     int            `json:"amount"`
	Status     Status         `json:"status"`
	Reference  string         `json:"reference"`
	Target     string         `json:"target"` // nomor tujuan topup / pulsa
	CreatedAt  time.Time      `bun:",nullzero,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt  time.Time      `bun:",nullzero,notnull,default:current_timestamp" json:"updated_at"`
	PaidAt     time.Time      `json:"paid_at"`
	CompletedAt time.Time     `json:"completed_at"`
}

type CreateInput struct {
	ProductID uuid.UUID        `form:"product_id" json:"product_id" validate:"required"`
	Channel   payment.Channel  `form:"channel" json:"channel" validate:"required,oneof=qris va ewallet"`
	Target    string           `form:"target" json:"target" validate:"required"`
}

type Service struct {
	db           *bun.DB
	validate     *validator.Validate
	payments     *payment.Service
	rewards      *reward.Service
	users        *user.Service
	productStore *product.Service
}

func NewService(db *bun.DB, validate *validator.Validate, payments *payment.Service, rewards *reward.Service, users *user.Service, products *product.Service) *Service {
	return &Service{
		db:           db,
		validate:     validate,
		payments:     payments,
		rewards:      rewards,
		users:        users,
		productStore: products,
	}
}

func (s *Service) Create(ctx context.Context, userID uuid.UUID, in CreateInput) (*Order, *payment.Invoice, error) {
	if err := s.validate.Struct(&in); err != nil {
		return nil, nil, err
	}

	product, err := s.productStore.FindByID(ctx, in.ProductID)
	if err != nil {
		return nil, nil, err
	}

	order := &Order{
		ID:        uuid.New(),
		UserID:    userID,
		ProductID: product.ID,
		Channel:   in.Channel,
		Amount:    product.Price,
		Status:    StatusPending,
		Target:    in.Target,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if _, err := s.db.NewInsert().Model(order).Exec(ctx); err != nil {
		return nil, nil, err
	}

	invoice, err := s.payments.CreateInvoice(ctx, order.ID, order.Amount, order.Channel)
	if err != nil {
		return order, nil, err
	}

	return order, invoice, nil
}

func (s *Service) MarkPaid(ctx context.Context, orderID uuid.UUID) error {
	var ord Order
	if err := s.db.NewSelect().Model(&ord).Where("id = ?", orderID).Scan(ctx); err != nil {
		return err
	}

	if _, err := s.db.NewUpdate().
		Model((*Order)(nil)).
		Set("status = ?", StatusPaid).
		Set("paid_at = ?", time.Now()).
		Set("updated_at = ?", time.Now()).
		Where("id = ?", orderID).
		Exec(ctx); err != nil {
		return err
	}

	// Trigger reward accrual asynchronously; errors are logged inside rewards.
	go s.rewards.HandleOrderPaid(context.Background(), ord.UserID, ord.Amount)
	return nil
}

func (s *Service) ListByUser(ctx context.Context, userID uuid.UUID) ([]Order, error) {
	var orders []Order
	err := s.db.NewSelect().
		Model(&orders).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Scan(ctx)
	return orders, err
}
