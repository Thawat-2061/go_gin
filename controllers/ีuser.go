package controllers

import (
	"go-gin/handlers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func UserController(router *gin.Engine, db *gorm.DB) {
	handler := handlers.NewUserHandler(db)

	userGroup := router.Group("/users")
	{
		userGroup.POST("/register", handler.Register)
		userGroup.POST("/login", handler.Login) // ต้องมีบรรทัดนี้
	}
}
