package routes

import (
	controller "github.com/1107-adishjain/golang-jwt/internal/controllers"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AuthRoutes(IncomingRoutes *gin.Engine, db *gorm.DB) {
	IncomingRoutes.POST("/api/signup", controller.SignUp(db))
	IncomingRoutes.POST("/api/login", controller.Login(db))
	IncomingRoutes.GET("/api/google-login", controller.GoogleLogin(db))
	IncomingRoutes.GET("/google/callback", controller.GoogleCallback(db))
}
