package userhttp

import (
	"context"
	"github.com/Zhoangp/Auth-Service/pb"
)

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
