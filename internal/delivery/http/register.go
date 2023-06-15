package userhttp

import (
	"context"
	"fmt"
	"github.com/Zhoangp/Auth-Service/internal/model"
	"github.com/Zhoangp/Auth-Service/pb"
	"github.com/Zhoangp/Auth-Service/pb/cart"
	"github.com/Zhoangp/Auth-Service/pb/mail"
	"github.com/Zhoangp/Auth-Service/pkg/client"
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
	cartService, err := client.InitCartServiceClient(userHandler.cf)
	resCreateCart, err := cartService.CreateCart(ctx, &cart.CreateCartRequest{
		UserId: data.FakeId,
	})
	fmt.Println(resCreateCart)
	if err != nil {
		fmt.Println(err)
	}
	mailService, err := client.InitServiceClient(userHandler.cf)
	if err != nil {
		fmt.Println(err)
		return &pb.RegisterResponse{
			Error: HandleError(err),
		}, nil
	}
	res, err := mailService.SendTokenVerifyAccount(ctx, &mail.SendTokenVerifyAccountRequest{
		Mail: &mail.Mail{
			DestMail: data.Email,
			Subject:  "Verify Account",
		},
		Token: token,
		Name:  data.LastName + " " + data.FirstName,
		Url:   "http://" + userHandler.cf.ClientSide.URL + "/register/successverify?token=",
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
