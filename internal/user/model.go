package user

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Role string

const (
	RoleDriver Role = "driver"
	RoleMitra  Role = "mitra"
	RoleAdmin  Role = "admin"
)

// User represents a driver or mitigation account in the system.
type User struct {
	bun.BaseModel `bun:"table:users"`

	ID           uuid.UUID `bun:",pk,type:uuid,default:gen_random_uuid()" json:"id"`
	Name         string    `json:"name"`
	Phone        string    `json:"phone"`
	VehicleType  string    `json:"vehicle_type"`
	PasswordHash string    `json:"-"`
	Role         Role      `json:"role"`
	Points       int       `json:"points"`
	CreatedAt    time.Time `bun:",nullzero,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt    time.Time `bun:",nullzero,notnull,default:current_timestamp" json:"updated_at"`
}

type PublicUser struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Phone       string    `json:"phone"`
	VehicleType string    `json:"vehicle_type"`
	Role        Role      `json:"role"`
	Points      int       `json:"points"`
}

func (u User) Public() PublicUser {
	return PublicUser{
		ID:          u.ID,
		Name:        u.Name,
		Phone:       u.Phone,
		VehicleType: u.VehicleType,
		Role:        u.Role,
		Points:      u.Points,
	}
}
