package api

import (
	"database/sql"
	"ecoswap/api/handler"
	_ "ecoswap/api/docs"
	"github.com/gin-gonic/gin"
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
func Router(db *sql.DB) {
	router := gin.Default()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	h := handler.NewHandlerRepo(db)

	router.POST("/register", h.Register)
	router.POST("/login", h.Login)
}
