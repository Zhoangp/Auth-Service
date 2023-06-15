package userhttp

import (
	"github.com/Zhoangp/Auth-Service/config"
	"github.com/Zhoangp/Auth-Service/internal/model"
	"github.com/Zhoangp/Auth-Service/pb"
	"github.com/Zhoangp/Auth-Service/pb/err"
	"github.com/Zhoangp/Auth-Service/pkg/common"
	"github.com/Zhoangp/Auth-Service/pkg/utils"
)

type UserHandler struct {
	UC UserUseCase
	pb.UnimplementedAuthServiceServer
	cf *config.Config
}
type UserUseCase interface {
	Register(data *model.Users) (*model.Users, string, error)
	Login(data *model.UserLogin) (*utils.Token, *utils.Token, *model.Users, error)
	GetNewToken(refreshToken string) (*utils.Token, error)
	GetUserNotVerified(email string) error
	GetTokenVerify(email string, key string) (*model.Users, string, error)
}

func NewUserHandler(cf *config.Config, userUC UserUseCase) *UserHandler {
	return &UserHandler{cf: cf, UC: userUC}
}

func HandleError(e error) *err.ErrorResponse {
	if errors, ok := e.(*common.AppError); ok {
		return &err.ErrorResponse{
			Code:    int64(errors.StatusCode),
			Message: errors.Message,
		}
	}
	appErr := common.ErrInternal(e.(error))
	return &err.ErrorResponse{
		Code:    int64(appErr.StatusCode),
		Message: appErr.Message,
	}
}
