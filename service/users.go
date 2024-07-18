package service

import (
	"context"
	"database/sql"
	pb "ecoswap/genproto/users"
	"ecoswap/pkg/logger"
	"ecoswap/storage/postgres"
	r "ecoswap/storage/redis"
	"fmt"
	"log/slog"

	"github.com/redis/go-redis/v9"
)

type UsersService struct{
	pb.UnimplementedUsersServiceServer
	Logger *slog.Logger
	Redis *r.UserRedis
	User *postgres.UsersRepo
}

func NewUsersService(db *sql.DB, rdb *redis.Client)*UsersService{
	return &UsersService{
		Logger: logger.NewLogger(),
		Redis: r.NewRedisRepo(rdb),
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

func(U *UsersService) GetEcoPointsByUser(ctx context.Context, req *pb.UserId)(*pb.UserEcoPoints, error){
	resp, err := U.User.GetEcoPointsByUser(req)
	if err != nil{
		U.Logger.Error(fmt.Sprintf("Databazadan ma'lumotlarni olishda xato: %v", err))
		return nil, err
	}
	return resp, nil
}
func(U *UsersService) CreateEcoPointsByUser(ctx context.Context, req *pb.CreateEcoPoints)(*pb.InfoUserEcoPoints, error){
	fmt.Println(11)
	resp, err := U.User.CreateEcoPointsByUser(req)
	if err != nil{
		U.Logger.Error(fmt.Sprintf("Databazadan ma'lumotlarni olishda xato: %v", err))
		return nil, err
	}
	return resp, nil
}
func(U *UsersService) HistoryEcoPointsByUser(ctx context.Context, req *pb.HistoryReq)(*pb.Histories, error){
	resp, err := U.Redis.HistoryEcoPointsByUser(req)
	if err != nil{
		U.Logger.Error(fmt.Sprintf("Databazadan ma'lumotlarni olishda xato: %v", err))
		return nil, err
	}
	return resp, nil
}