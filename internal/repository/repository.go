package repository

import "gaspoll/internal/entity"

type UserRepository interface {
	ListUsers() ([]*entity.User, error)
	FindUserByID(id string) (*entity.User, error)
	FindUserByPhone(phone string) (*entity.User, error)
	Save(user *entity.User) error
}

type ProductRepository interface {
	ListProducts() ([]*entity.Product, error)
	FindProductByID(id string) (*entity.Product, error)
}

type OrderRepository interface {
	ListByUser(userID string) ([]*entity.Order, error)
	Save(order *entity.Order) error
}

type VoucherRepository interface {
	ListByUser(userID string) ([]*entity.Voucher, error)
	FindVoucherByID(id string) (*entity.Voucher, error)
	Save(voucher *entity.Voucher) error
}

type MitraRepository interface {
	ListMitras() ([]*entity.Mitra, error)
	FindMitraByID(id string) (*entity.Mitra, error)
}

type SettlementRepository interface {
	ListSettlements() ([]*entity.Settlement, error)
	Save(settlement *entity.Settlement) error
}

type Repository struct {
	User       UserRepository
	Product    ProductRepository
	Order      OrderRepository
	Voucher    VoucherRepository
	Mitra      MitraRepository
	Settlement SettlementRepository
}
