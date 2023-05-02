package main

import (
	"fmt"
	"github.com/Zhoangp/Auth-Service/config"
	userhttp "github.com/Zhoangp/Auth-Service/internal/delivery/http"
	"github.com/Zhoangp/Auth-Service/internal/repo"
	"github.com/Zhoangp/Auth-Service/internal/usecase"
	"github.com/Zhoangp/Auth-Service/pb"
	"github.com/Zhoangp/Auth-Service/pkg/client"
	"github.com/Zhoangp/Auth-Service/pkg/database/mysql"
	"github.com/Zhoangp/Auth-Service/pkg/utils"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

func main() {
	env := os.Getenv("ENV")
	fileName := "config/config-local.yml"
	if env == "app" {
		fileName = "config/config-app.yml"
	}

	cf, err := config.LoadConfig(fileName)
	if err != nil {
		log.Fatalln("Failed at config", err)
	}
	gormDb, err := mysql.NewMysql(cf)
	if err != nil {
		fmt.Println(err)
		return
	}
	lis, err := net.Listen("tcp", ":"+cf.Service.Port)
	if err != nil {
		fmt.Println(err)
		return
	}
	hasher := utils.NewHasher("courses", 3)
	fmt.Println("Auth Svc on", cf.Service.Port)

	clientMailService, err := client.InitServiceClient(cf)
	if err != nil {
		fmt.Println(err)
		return
	}
	repoUser := repo.NewUserRepository(gormDb)
	useCaseUser := usecase.NewUserUseCase(repoUser, cf, hasher)
	hdlUser := userhttp.NewUserHandler(cf, useCaseUser, clientMailService)

	grpcServer := grpc.NewServer()
	pb.RegisterAuthServiceServer(grpcServer, hdlUser)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalln("Failed to serve:", err)
	}
}
