package userhttp

import (
	"context"
	"fmt"
	"github.com/Zhoangp/Auth-Service/internal/model"
	"github.com/Zhoangp/Auth-Service/pb"
	"github.com/Zhoangp/Auth-Service/pb/mail"
)

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
	data, token, err := userHandler.UC.Register(&user)
	if err != nil {
		return &pb.RegisterResponse{
			Error: HandleError(err),
		}, nil
	}
	res, err := userHandler.mailClient.SendTokenVerifyAccount(ctx, &mail.SendTokenVerifyAccountRequest{
		Mail: &mail.Mail{
			DestMail: data.Email,
			Subject:  "Verify Account",
		},
		Token: token,
		Name:  data.LastName + " " + data.FirstName,
		Url:   "http://" + userHandler.cf.ClientSide.URL + "/courses/register/successverify?token=",
	})

	if err != nil {
		fmt.Println(err)
		return &pb.RegisterResponse{
			Error: HandleError(err),
		}, nil
	}
	if res.Error != nil {

		return &pb.RegisterResponse{
			Error: res.Error,
		}, nil
	}
	return &pb.RegisterResponse{}, nil
}
