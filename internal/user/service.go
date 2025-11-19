package user

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun"
	"golang.org/x/crypto/bcrypt"
)

// Service wraps user CRUD + authentication logic.
type Service struct {
	db       *bun.DB
	validate *validator.Validate
}

func NewService(db *bun.DB, validate *validator.Validate) *Service {
	return &Service{db: db, validate: validate}
}

type RegisterInput struct {
	Name        string `form:"name" json:"name" validate:"required"`
	Phone       string `form:"phone" json:"phone" validate:"required,min=10,max=16"`
	Password    string `form:"password" json:"password" validate:"required,min=6"`
	VehicleType string `form:"vehicle_type" json:"vehicle_type" validate:"required"`
	Role        Role   `form:"role" json:"role" validate:"required,oneof=driver mitra admin"`
}

type LoginInput struct {
	Phone    string `form:"phone" json:"phone" validate:"required"`
	Password string `form:"password" json:"password" validate:"required"`
}

func (s *Service) Register(ctx context.Context, in RegisterInput) (*User, error) {
	if err := s.validate.Struct(&in); err != nil {
		return nil, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &User{
		ID:           uuid.New(),
		Name:         strings.TrimSpace(in.Name),
		Phone:        strings.TrimSpace(in.Phone),
		PasswordHash: string(hash),
		VehicleType:  in.VehicleType,
		Role:         in.Role,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	_, err = s.db.NewInsert().Model(user).Exec(ctx)
	return user, err
}

func (s *Service) Authenticate(ctx context.Context, in LoginInput) (*User, error) {
	if err := s.validate.Struct(&in); err != nil {
		return nil, err
	}

	var user User
	err := s.db.NewSelect().
		Model(&user).
		Where("phone = ?", in.Phone).
		Limit(1).
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(in.Password)) != nil {
		return nil, errors.New("invalid credentials")
	}

	return &user, nil
}

func (s *Service) AddPoints(ctx context.Context, userID uuid.UUID, amount int) error {
	_, err := s.db.NewUpdate().
		Model((*User)(nil)).
		Set("points = points + ?", amount).
		Set("updated_at = ?", time.Now()).
		Where("id = ?", userID).
		Exec(ctx)
	if err != nil {
		log.Error().Err(err).Msg("failed to add points")
	}
	return err
}

func (s *Service) FindByID(ctx context.Context, userID uuid.UUID) (*User, error) {
	var user User
	err := s.db.NewSelect().Model(&user).Where("id = ?", userID).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
