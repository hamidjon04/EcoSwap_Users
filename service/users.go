package service

import (
	"context"
	"database/sql"
	pb "ecoswap/genproto/users"
	"ecoswap/pkg/logger"
	"ecoswap/storage/postgres"
	"fmt"
	"log/slog"
)

type UsersService struct{
	pb.UnimplementedUsersServiceServer
	Logger *slog.Logger
	User *postgres.UsersRepo
}

func NewUsersService(db *sql.DB)*UsersService{
	return &UsersService{
		Logger: logger.NewLogger(),
		User: postgres.NewUsersRepo(db),
	}
}

func(U *UsersService) GetProfile(ctx context.Context, req *pb.UserId)(*pb.UserInfo, error){
	resp, err := U.User.GetProfile(req)
	if err != nil{
		U.Logger.Error(fmt.Sprintf("Databazadan ma'lumotlarni olishda xato: %v", err))
		return nil, err
	}
	return resp, nil
}

func(U *UsersService) UpdateProfile(ctx context.Context, req *pb.ProfileUpdate)(*pb.UpdateResponse, error){
	resp, err := U.User.UpdateProfile(req)
	if err != nil{
		U.Logger.Error(fmt.Sprintf("Databazadan ma'lumotlarni olishda xato: %v", err))
		return nil, err
	}
	return resp, nil
}

func(U *UsersService) DeleteProfile(ctx context.Context, req *pb.UserId)(*pb.Status, error){
	resp, err := U.User.DeleteProfile(req)
	if err != nil{
		U.Logger.Error(fmt.Sprintf("Databazadan ma'lumotlarni olishda xato: %v", err))
		return nil, err
	}
	return resp, nil
}

func(U *UsersService) GetAllUsers(ctx context.Context, req *pb.FilterField)(*pb.Users, error){
	resp, err := U.User.GetAllUsers(req)
	if err != nil{
		U.Logger.Error(fmt.Sprintf("Databazadan ma'lumotlarni olishda xato: %v", err))
		return nil, err
	}
	return resp, nil
}
