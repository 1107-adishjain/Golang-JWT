package helpers

import (
	jwt "github.com/golang-jwt/jwt/v5"
	"os"
	"time"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

type Claims struct {
	UserID    string `json:"user_id"`
	UserType  string `json:"user_type"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	jwt.RegisteredClaims
}

func GenerateJWT(userid, usertype, firstname, lastname string) (access_token, refresh_token string, err error) {
	accessExp := time.Now().Add(15 * time.Minute)
	refreshExp := time.Now().Add(7 * 24 * time.Hour)

	accessClaims := Claims{
		UserID:    userid,
		UserType:  usertype,
		FirstName: firstname,
		LastName:  lastname,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessExp),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   userid,
		},
	}

	refreshClaims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(refreshExp),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Subject:   userid,
	}

	accessTokenJwt := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	refreshTokenJwt := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)

	access_token, err = accessTokenJwt.SignedString(jwtSecret)
	if err != nil {
		return "", "", err
	}
	refresh_token, err = refreshTokenJwt.SignedString(jwtSecret)
	if err != nil {
		return "", "", err
	}
	return access_token, refresh_token, nil
}
