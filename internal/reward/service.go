package reward

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun"

	"gaspoll/internal/config"
	"gaspoll/internal/user"
	"gaspoll/internal/voucher"
)

type Progress struct {
	TotalSpend   int
	NextVoucher  int
	Points       int
	RecentReward *voucher.Voucher
}

type Service struct {
	db       *bun.DB
	cfg      config.Config
	users    *user.Service
	vouchers *voucher.Service
}

func NewService(db *bun.DB, cfg config.Config, users *user.Service, vouchers *voucher.Service) *Service {
	return &Service{db: db, cfg: cfg, users: users, vouchers: vouchers}
}

// HandleOrderPaid increments points and issues voucher when threshold reached.
func (s *Service) HandleOrderPaid(ctx context.Context, userID uuid.UUID, orderAmount int) {
	if err := s.users.AddPoints(ctx, userID, orderAmount/1000); err != nil {
		log.Error().Err(err).Msg("failed to add points for paid order")
	}

	total, err := s.totalSpend(ctx, userID)
	if err != nil {
		log.Error().Err(err).Msg("failed to calculate total spend")
		return
	}

	if total >= s.cfg.RewardTarget {
		v, err := s.vouchers.Issue(ctx, userID, s.cfg.RewardValue, "Bengkel Mitra Gaspoll")
		if err != nil {
			log.Error().Err(err).Msg("failed to issue voucher")
			return
		}
		log.Info().Str("voucher_code", v.Code).Msg("voucher issued")
	}
}

func (s *Service) totalSpend(ctx context.Context, userID uuid.UUID) (int, error) {
	var total int
	err := s.db.NewSelect().
		Table("orders").
		ColumnExpr("COALESCE(sum(amount),0)").
		Where("user_id = ?", userID).
		Scan(ctx, &total)
	return total, err
}

func (s *Service) Progress(ctx context.Context, userID uuid.UUID) (Progress, error) {
	total, err := s.totalSpend(ctx, userID)
	if err != nil {
		return Progress{}, err
	}

	progress := Progress{
		TotalSpend:  total,
		NextVoucher: s.cfg.RewardTarget - (total % s.cfg.RewardTarget),
	}

	usr, err := s.users.FindByID(ctx, userID)
	if err == nil {
		progress.Points = usr.Points
	}

	vouchers, _ := s.vouchers.ListByUser(ctx, userID)
	if len(vouchers) > 0 {
		progress.RecentReward = &vouchers[0]
	}
	return progress, nil
}

// ExpireOverdue should run in background (cron) to set expired vouchers.
func (s *Service) ExpireOverdue(ctx context.Context) error {
	_, err := s.db.NewUpdate().
		Model((*voucher.Voucher)(nil)).
		Set("status = ?", voucher.StatusExpired).
		Where("expires_at < ? AND status = ?", time.Now(), voucher.StatusIssued).
		Exec(ctx)
	return err
}
