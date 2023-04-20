package userhttp

import (
	"context"
	"github.com/Zhoangp/Auth-Service/internal/model"
	"github.com/Zhoangp/Auth-Service/pb"
	"github.com/Zhoangp/Auth-Service/pkg/common"
	"github.com/Zhoangp/Auth-Service/pkg/utils"
	"time"
)

type UserHandler struct {
	UC UserUseCase
	pb.UnimplementedAuthServiceServer
}
type UserUseCase interface {
	Register(data *model.Users) error
	Login(data *model.UserLogin) (*utils.Token, *utils.Token, *model.Users, error)
	GetNewToken(refreshToken string) (*utils.Token, error)
	GetUserNotVerified(email string) error
	GetTokenVerify(email string) error
}

func NewUserHandler(userUC UserUseCase) *UserHandler {
	return &UserHandler{UC: userUC}
}

func (userHandler *UserHandler) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	user := model.Users{
		Email:     req.Email,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Phone:     req.PhoneNumber,
		Role:      req.Role,
		Address:   req.Address,
	}
	if err := userHandler.UC.Register(&user); err != nil {
		return &pb.RegisterResponse{
			Error: HandleError(err),
		}, nil
	}
	return &pb.RegisterResponse{}, nil
}
func (userHandler *UserHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	var data model.UserLogin

	data.Email = req.Email
	data.Password = req.Password
	token, refreshToken, user, err := userHandler.UC.Login(&data)
	if err != nil {
		return &pb.LoginResponse{
			Error: HandleError(err),
		}, nil
	}
	return &pb.LoginResponse{
		Information: &pb.User{
			FirstName:   user.FirstName,
			LastName:    user.LastName,
			Email:       user.Email,
			PhoneNumber: user.Password,
			Address:     user.Password,
			Avatar: &pb.Picture{
				Url:    user.Avatar.Url,
				Width:  user.Avatar.Width,
				Height: user.Avatar.Height,
			},
			LastLogin: user.LastLogin.Format(time.Stamp),
		},
		AccessToken:  token.AccessToken,
		RefreshToken: refreshToken.AccessToken,
		ExpiresAt:    uint32(token.ExpiresAt),
		TokenType:    "Bearer",
	}, nil
}

func (userHandler *UserHandler) NewToken(ctx context.Context, req *pb.NewTokenRequest) (*pb.NewTokenResponse, error) {

	token, err := userHandler.UC.GetNewToken(req.RefreshToken)
	if err != nil {
		return &pb.NewTokenResponse{
			Error: HandleError(err),
		}, nil
	}
	return &pb.NewTokenResponse{
		AccessToken: token.AccessToken,
		ExpiresAt:   uint32(token.ExpiresAt),
	}, nil
}

func (userHandler *UserHandler) VerifyAccount(ctx context.Context, req *pb.VerifyAccountRequest) (*pb.VerifyAccountResponse, error) {
	if err := userHandler.UC.GetUserNotVerified(req.Email); err != nil {
		return &pb.VerifyAccountResponse{
			Error: HandleError(err),
		}, nil
	}
	return &pb.VerifyAccountResponse{}, nil
}
func (userHandler *UserHandler) GetTokenVeriryAccount(ctx context.Context, req *pb.VerifyAccountRequest) (*pb.VerifyAccountResponse, error) {
	if err := userHandler.UC.GetTokenVerify(req.Email); err != nil {
		return &pb.VerifyAccountResponse{
			Error: HandleError(err),
		}, nil
	}
	return &pb.VerifyAccountResponse{}, nil

}
func HandleError(err error) *pb.ErrorResponse {
	if errors, ok := err.(*common.AppError); ok {
		return &pb.ErrorResponse{
			Code:    int64(errors.StatusCode),
			Message: errors.Message,
		}
	}
	appErr := common.ErrInternal(err.(error))
	return &pb.ErrorResponse{
		Code:    int64(appErr.StatusCode),
		Message: appErr.Message,
	}
}
