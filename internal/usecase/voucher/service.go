package voucher

import (
	"errors"
	"time"

	"gaspoll/internal/entity"
	"gaspoll/internal/repository"
)

type Service struct {
	Vouchers repository.VoucherRepository
}

func (s Service) ListByUser(userID string) ([]*entity.Voucher, error) {
	return s.Vouchers.ListByUser(userID)
}

func (s Service) Redeem(id string) (*entity.Voucher, error) {
	v, err := s.Vouchers.FindVoucherByID(id)
	if err != nil {
		return nil, err
	}
	if v.IsRedeemed {
		return nil, errors.New("voucher sudah digunakan")
	}
	now := time.Now()
	v.IsRedeemed = true
	v.RedeemedAt = &now
	if err := s.Vouchers.Save(v); err != nil {
		return nil, err
	}
	return v, nil
}
