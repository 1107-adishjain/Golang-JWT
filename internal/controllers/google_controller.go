package controllers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GoogleLogin(db *gorm.DB) gin.HandlerFunc{
	return func(c *gin.Context){
		// this function would handle the google login logic
	}
}