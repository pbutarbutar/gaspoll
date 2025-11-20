package settlement

import (
	"gaspoll/internal/entity"
	"gaspoll/internal/repository"
)

type Service struct {
	Settlements repository.SettlementRepository
}

func (s Service) List() ([]*entity.Settlement, error) {
	return s.Settlements.ListSettlements()
}

func (s Service) Add(settlement *entity.Settlement) error {
	return s.Settlements.Save(settlement)
}
