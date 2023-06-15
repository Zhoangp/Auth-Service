package client

import (
	"fmt"
	"github.com/Zhoangp/Auth-Service/config"
	"github.com/Zhoangp/Auth-Service/pb/cart"
	"github.com/Zhoangp/Auth-Service/pb/mail"

	"google.golang.org/grpc"
)

func InitServiceClient(c *config.Config) (mail.MailServiceClient, error) {
	// using WithInsecure() because no SSL running
	cc, err := grpc.Dial(c.OtherServices.MailServiceURL, grpc.WithInsecure())
	if err != nil {
		fmt.Println("Could not connect:", err)
		return nil, err
	}
	return mail.NewMailServiceClient(cc), nil
}
func InitCartServiceClient(c *config.Config) (cart.CartServiceClient, error) {
	// using WithInsecure() because no SSL running
	cc, err := grpc.Dial(c.OtherServices.CartServiceUrl, grpc.WithInsecure())
	if err != nil {
		fmt.Println("Could not connect:", err)
		return nil, err
	}
	fmt.Println(cc)
	return cart.NewCartServiceClient(cc), nil
}
