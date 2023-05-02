package usecase

import (
	"github.com/Zhoangp/Auth-Service/internal/model"
	"github.com/Zhoangp/Auth-Service/pkg/common"
	"github.com/Zhoangp/Auth-Service/pkg/utils"
	"time"
)

func (uc *userUseCase) Register(data *model.Users) (*model.Users, string, error) {
	if user, _ := uc.userRepo.FindDataWithCondition(map[string]any{"email": data.Email}); user != nil {
		return nil, "", model.ErrEmailExisted
	}
	if err := data.PrepareCreate(); err != nil {
		return nil, "", err
	}
	data.Avatar = &common.Image{
		Id:     1,
		Url:    "https://courses-storages.s3.ap-northeast-1.amazonaws.com/default-img/images.png",
		Width:  "250px",
		Height: "250px",
	}
	data.LastLogin = time.Now()
	if err := uc.userRepo.NewUsers(data); err != nil {
		return nil, "", err
	}
	token, err := utils.GenerateToken(utils.TokenPayload{Email: data.Email, Role: data.Role, Password: data.Password, Verified: false}, uc.cf.Service.ActiveTokenExpired, uc.cf.Service.Secret)
	if err != nil {
		return nil, "", err
	}
	//if err := utils.SendToken(uc.cf, data.Email, token.AccessToken, data.FirstName+data.LastName, "http://127.0.0.1:8080/courses/register/successverify?token="); err != nil {
	//	return nil, "", err
	//}

	return data, token.AccessToken, nil
}
func (uc *userUseCase) GetTokenVerify(email string, key string) (*model.Users, string, error) {
	user, err := uc.userRepo.FindDataWithCondition(map[string]any{"email": email})
	if err != nil {
		return nil, "", model.ErrEmailOrPasswordInvalid
	}
	payload := utils.TokenPayload{Email: user.Email, Role: user.Role, Password: user.Password, Verified: user.Verified, Key: key}
	token, err := utils.GenerateToken(payload, uc.cf.Service.ActiveTokenExpired, uc.cf.Service.Secret)
	if err != nil {
		return nil, "", err
	}

	return user, token.AccessToken, nil
}
