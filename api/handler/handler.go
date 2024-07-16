package handler

import (
	"database/sql"
	"ecoswap/pkg/logger"
	"ecoswap/storage/postgres"
	pb "ecoswap/storage/redis"
	"log/slog"

	"github.com/redis/go-redis/v9"
)

type Handler struct {
	UserRepo *postgres.UsersRepo
	Redis    *pb.UserRedis
	Logger   *slog.Logger
}

func NewHandlerRepo(db *sql.DB, rdb *redis.Client) *Handler {
	return &Handler{
		UserRepo: postgres.NewUsersRepo(db),
		Redis:    pb.NewRedisRepo(rdb),
		Logger:   logger.NewLogger(),
	}
}
