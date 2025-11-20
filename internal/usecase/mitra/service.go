package mitra

import (
	"gaspoll/internal/entity"
	"gaspoll/internal/repository"
)

type Service struct {
	Mitras repository.MitraRepository
}

func (s Service) List() ([]*entity.Mitra, error) {
	return s.Mitras.ListMitras()
}

func (s Service) Find(id string) (*entity.Mitra, error) {
	return s.Mitras.FindMitraByID(id)
}
