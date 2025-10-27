package routes

import (
	controller "github.com/1107-adishjain/golang-jwt/internal/controllers"
	"github.com/1107-adishjain/golang-jwt/internal/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func UserRoutes(IncomingRoutes *gin.Engine, db *gorm.DB) {
	IncomingRoutes.Use(middleware.Authenticate()) ///the authMiddlware will be used to check whther the user creating request on the user and user id get is an verfified user and this will be verified using the authetnicate() fucntion that is created once the user verfies then it will allow the user to access these routes nor they will be blocked.
	IncomingRoutes.GET("/api/users", controller.GetUsers(db))
	IncomingRoutes.GET("/api/user/:id", controller.GetUser(db))
}
