package server

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun"

	"gaspoll/internal/auth"
	"gaspoll/internal/config"
	"gaspoll/internal/mitra"
	"gaspoll/internal/order"
	"gaspoll/internal/payment"
	"gaspoll/internal/product"
	"gaspoll/internal/reward"
	"gaspoll/internal/settlement"
	"gaspoll/internal/user"
	"gaspoll/internal/view"
	"gaspoll/internal/voucher"
)

type requestValidator struct {
	validate *validator.Validate
}

func (rv *requestValidator) Validate(i interface{}) error {
	return rv.validate.Struct(i)
}

// New bootstraps Echo with middleware, renderer, and routes.
func New(cfg config.Config, db *bun.DB) (*echo.Echo, error) {
	e := echo.New()
	e.HideBanner = true
	e.Validator = &requestValidator{validate: validator.New()}

	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	e.Use(middleware.Secure())

	sessionManager := auth.NewSessionManager(cfg.SessionSecret)
	e.Use(sessionManager.Middleware)

	renderer, err := view.NewEngine()
	if err != nil {
		return nil, err
	}
	e.Renderer = renderer

	// services
	users := user.NewService(db, validator.New())
	vouchers := voucher.NewService(db)
	rewards := reward.NewService(db, cfg, users, vouchers)
	products := product.NewService(db)
	payments := payment.NewService(cfg.PaymentBaseURL)
	orders := order.NewService(db, validator.New(), payments, rewards, users, products)
	mitras := mitra.NewService(db)
	settlements := settlement.NewService(db)

	svcs := services{
		users:       users,
		vouchers:    vouchers,
		rewards:     rewards,
		products:    products,
		payments:    payments,
		orders:      orders,
		mitras:      mitras,
		settlements: settlements,
		session:     sessionManager,
	}

	registerRoutes(e, renderer, svcs)
	log.Info().Msg("server initialized")

	return e, nil
}
