package userhttp

import (
	"context"
	"fmt"
	"github.com/Zhoangp/Auth-Service/pb"
	"github.com/Zhoangp/Auth-Service/pb/mail"
)

func (userHandler *UserHandler) GetTokenResetPassword(ctx context.Context, req *pb.VerifyAccountRequest) (*pb.VerifyAccountResponse, error) {
	user, token, err := userHandler.UC.GetTokenVerify(req.Email, "forget")
	fmt.Println("abc")
	res, err := userHandler.mailClient.SendTokenVerifyAccount(ctx, &mail.SendTokenVerifyAccountRequest{
		Mail: &mail.Mail{
			DestMail: user.Email,
			Subject:  "Reset Password",
		},
		Token: token,
		Name:  user.FirstName + " " + user.LastName,
		Url:   "http://" + userHandler.cf.ClientSide.URL + "/courses/recover?token=",
	})

	if err != nil {
		return &pb.VerifyAccountResponse{
			Error: HandleError(err),
		}, nil
	}
	if res.Error != nil {
		return &pb.VerifyAccountResponse{
			Error: res.Error,
		}, nil
	}
	return &pb.VerifyAccountResponse{}, nil
}
