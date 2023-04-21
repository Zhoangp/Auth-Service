package usecase

import (
	"github.com/Zhoangp/Auth-Service/config"
	"github.com/Zhoangp/Auth-Service/internal/model"
	"github.com/Zhoangp/Auth-Service/internal/repo"
	"github.com/Zhoangp/Auth-Service/pkg/common"
	"github.com/Zhoangp/Auth-Service/pkg/utils"
	"time"
)

type UserRepository interface {
	NewUsers(data *model.Users) error
	FindDataWithCondition(conditions map[string]any) (*model.Users, error)
	UpdateUser(user model.Users, newInformation map[string]any) error
}

type userUseCase struct {
	cf       *config.Config
	userRepo *repo.UserRepository
}

func NewUserUseCase(userRepo *repo.UserRepository, cf *config.Config) *userUseCase {
	return &userUseCase{cf, userRepo}
}
func (uc *userUseCase) Register(data *model.Users) error {
	if user, _ := uc.userRepo.FindDataWithCondition(map[string]any{"email": data.Email}); user != nil {
		return model.ErrEmailExisted
	}
	if err := data.PrepareCreate(); err != nil {
		return err
	}
	data.Avatar = &common.Image{
		Id:     1,
		Url:    "https://courses-storages.s3.ap-northeast-1.amazonaws.com/default-img/images.png",
		Width:  "250px",
		Height: "250px",
	}
	data.LastLogin = time.Now()
	if err := uc.userRepo.NewUsers(data); err != nil {
		return err
	}
	token, err := utils.GenerateToken(utils.TokenPayload{Email: data.Email, Role: data.Role, Password: data.Password, Verified: false}, uc.cf.Service.ActiveTokenExpired, uc.cf.Service.Secret)
	if err != nil {
		return err
	}
	if err := utils.SendToken(uc.cf, data.Email, token.AccessToken, data.FirstName+data.LastName); err != nil {
		return err
	}

	return nil
}
func (uc *userUseCase) Login(data *model.UserLogin) (*utils.Token, *utils.Token, *model.Users, error) {
	user, err := uc.userRepo.FindDataWithCondition(map[string]any{"email": data.Email})
	if err != nil {
		return nil, nil, nil, model.ErrEmailOrPasswordInvalid
	}
	if !user.Verified {
		return nil, nil, nil, model.ErrAccountNotVerified
	}

	if err := utils.ComparePassword(user.Password, data.Password); err != nil {
		return nil, nil, nil, model.ErrEmailOrPasswordInvalid
	}
	lastLogin := user.LastLogin
	token, err := utils.GenerateToken(utils.TokenPayload{Email: user.Email, Role: user.Role, Password: user.Password}, uc.cf.Service.AccessTokenExpiredIn, uc.cf.Service.Secret)
	if err != nil {
		return nil, nil, nil, common.ErrInternal(err)
	}
	refreshToken, err := utils.GenerateToken(utils.TokenPayload{Email: user.Email, Role: user.Role, Password: user.Password}, uc.cf.Service.RefreshTokenExpiredIn, uc.cf.Service.Secret)
	if err != nil {
		return nil, nil, nil, common.ErrInternal(err)
	}
	if err := uc.userRepo.UpdateUser(user, map[string]any{"lastLogin": time.Now()}); err != nil {
		return nil, nil, nil, common.ErrInternal(err)
	}
	user.LastLogin = lastLogin
	return token, refreshToken, user, nil
}
func (uc *userUseCase) GetTokenVerify(email string) error {
	user, err := uc.userRepo.FindDataWithCondition(map[string]any{"email": email})
	if err != nil {
		return model.ErrEmailOrPasswordInvalid
	}
	payload := utils.TokenPayload{Email: user.Email, Role: user.Role, Password: user.Password, Verified: user.Verified}
	token, err := utils.GenerateToken(payload, uc.cf.Service.ActiveTokenExpired, uc.cf.Service.Secret)
	if err != nil {
		return err
	}
	if err := utils.SendToken(uc.cf, user.Email, token.AccessToken, user.LastName); err != nil {
		return err
	}
	return nil
}
func (uc *userUseCase) GetNewToken(refreshToken string) (*utils.Token, error) {
	data, err := utils.ValidateJWT(refreshToken, uc.cf)
	if err != nil {
		return nil, err
	}
	user, err := uc.userRepo.FindDataWithCondition(map[string]any{"email": data.Email})
	if err != nil {
		return nil, err
	}
	token, err := utils.GenerateToken(utils.TokenPayload{Email: user.Email, Role: user.Role}, uc.cf.Service.AccessTokenExpiredIn, uc.cf.Service.Secret)
	if err != nil {
		return nil, err
	}
	return token, err
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
