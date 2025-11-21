package entity

import "time"

type User struct {
	ID           string
	Name         string
	Phone        string
	VehicleType  string
	Points       int
	RewardGoal   int
	RewardStarts int
	CreatedAt    time.Time
}

type Product struct {
	ID          string
	Name        string
	Category    string
	Provider    string
	Price       int
	Description string
	Tags        []string
}

type Order struct {
	ID             string
	UserID         string
	ProductID      string
	ProductName    string
	Amount         int
	PaymentMethod  string
	Status         string
	WhatsAppLink   string
	PaymentCode    string
	PaymentPayload string
	RewardEarned   int
	CreatedAt      time.Time
}

type Voucher struct {
	ID          string
	Code        string
	UserID      string
	MitraID     string
	Title       string
	Description string
	PointsUsed  int
	IsRedeemed  bool
	RedeemedAt  *time.Time
	ExpiredAt   time.Time
	QRData      string
}

type RewardProgress struct {
	CurrentPoints int
	Goal          int
	NextReward    string
}

type Mitra struct {
	ID              string
	Name            string
	Address         string
	Location        string
	Promo           string
	Services        []string
	Distance        string
	BasePrice       int    // typical price for main service in IDR
	DiscountPercent int    // discount percent applied by promo
	PromoDetail     string // extra promo description / conditions
	Packages        []Package
}

type Package struct {
	ID              string
	Name            string
	Description     string
	BasePrice       int
	DiscountPercent int
}

type Settlement struct {
	ID           string
	MitraID      string
	Period       string
	TotalVoucher int
	Amount       int
	Status       string
}

type PaymentInstruction struct {
	Method   string
	Code     string
	DeepLink string
	Steps    []string
	QRString string
}
