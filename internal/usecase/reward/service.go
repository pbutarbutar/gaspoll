package reward

import "gaspoll/internal/entity"

type Service struct{}

func (Service) AddPoints(user *entity.User, earned int) {
	user.Points += earned
}

func (Service) Progress(user *entity.User) entity.RewardProgress {
	return entity.RewardProgress{
		CurrentPoints: user.Points,
		Goal:          user.RewardGoal,
		NextReward:    "Voucher servis bengkel mitra",
	}
}

func (Service) ShouldIssueVoucher(user *entity.User) bool {
	return user.Points >= user.RewardGoal
}

func (Service) ConsumeForVoucher(user *entity.User) {
	if user.Points >= user.RewardGoal {
		user.Points -= user.RewardGoal
	}
}
