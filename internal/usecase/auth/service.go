package auth

import (
	"fmt"
	"strings"
	"time"

	"gaspoll/internal/constant"
	"gaspoll/internal/entity"
	"gaspoll/internal/repository"
)

type Service struct {
	Users repository.UserRepository
}

func (s Service) Login(phone string) (*entity.User, error) {
	normalized := normalizePhone(phone)
	return s.Users.FindUserByPhone(normalized)
}

func (s Service) Register(name, phone, vehicle string) (*entity.User, error) {
	normalized := normalizePhone(phone)
	user := &entity.User{
		ID:          fmt.Sprintf("user-%d", time.Now().UnixNano()),
		Name:        strings.TrimSpace(name),
		Phone:       normalized,
		VehicleType: vehicle,
		Points:      0,
		RewardGoal:  constant.DefaultRewardGoal,
		CreatedAt:   time.Now(),
	}
	if err := s.Users.Save(user); err != nil {
		return nil, err
	}
	return user, nil
}

func normalizePhone(phone string) string {
	p := strings.TrimSpace(phone)
	p = strings.TrimPrefix(p, "+62")
	p = strings.TrimPrefix(p, "62")
	if strings.HasPrefix(p, "0") {
		return p
	}
	return "0" + p
}
