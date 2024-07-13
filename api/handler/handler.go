package handler

import (
	"database/sql"
	"ecoswap/pkg/logger"
	"ecoswap/storage/postgres"
	"log/slog"
)

type Handler struct {
	UserRepo postgres.UsersRepo
	Logger   *slog.Logger
}

func NewHandlerRepo(db *sql.DB) *Handler {
	return &Handler{
		UserRepo: *postgres.NewUsersRepo(db),
		Logger: logger.NewLogger(),
	}
}
