package api

import (
	"database/sql"
	_ "ecoswap/api/docs"
	"ecoswap/api/handler"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Auth Service
// @version 1.0
// @description This is the Auth service of EcoSwap project

// @contact.name Hamidjon
// @contact.email nuriddinovhamidjon2@gmail.com

// @host localhost:7777
// @BasePath /users
func Router(db *sql.DB, rdb *redis.Client) *gin.Engine{
	router := gin.Default()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	h := handler.NewHandlerRepo(db, rdb)

	router.POST("/register", h.Register)
	router.POST("/login", h.Login)

	return router
}
