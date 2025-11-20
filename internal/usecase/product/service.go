package product

import (
	"gaspoll/internal/entity"
	"gaspoll/internal/repository"
)

type Service struct {
	Products repository.ProductRepository
}

func (s Service) List() ([]*entity.Product, error) {
	return s.Products.ListProducts()
}

func (s Service) Find(id string) (*entity.Product, error) {
	return s.Products.FindProductByID(id)
}
