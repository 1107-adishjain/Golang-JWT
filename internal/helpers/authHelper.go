package helpers

import (
	"errors"
	// "go/token"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func ValidateUserId(c *gin.Context, userId string) (err error) {
	userType := c.GetString("user_type")
	uuid := c.GetString("user_id")
	if uuid != userId {
		return errors.New("you can only access your own user data")
	}
	if userType != "ADMIN" {
		return errors.New("only admin users can access this resource")
	}
	if userId == "" {
		return errors.New("user ID cannot be empty")
	}
	return nil
}

func VerifyJWT(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}
	return claims, nil
}
