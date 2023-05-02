package usecase

import (
	"github.com/Zhoangp/Auth-Service/config"
	"github.com/Zhoangp/Auth-Service/internal/model"
	"github.com/Zhoangp/Auth-Service/internal/repo"
	"github.com/Zhoangp/Auth-Service/pkg/common"
	"github.com/Zhoangp/Auth-Service/pkg/utils"
)

type UserRepository interface {
	NewUsers(data *model.Users) error
	FindDataWithCondition(conditions map[string]any) (*model.Users, error)
	UpdateUser(user model.Users, newInformation map[string]any) error
}

type userUseCase struct {
	cf       *config.Config
	userRepo *repo.UserRepository
	h        *utils.Hasher
}

func NewUserUseCase(userRepo *repo.UserRepository, cf *config.Config, h *utils.Hasher) *userUseCase {
	return &userUseCase{cf, userRepo, h}
}

func (uc *userUseCase) GetUserNotVerified(email string) error {
	user, err := uc.userRepo.FindDataWithCondition(map[string]any{"email": email, "verified": 0})
	if err != nil {
		return common.ErrEntityNotFound("Email", err)
	}
	if err := uc.userRepo.UpdateUser(user, map[string]any{"verified": 1}); err != nil {
		return common.ErrInternal(err)
	}
	return nil
}
