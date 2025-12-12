package controllers

import (
	"net/http"

	"github.com/1107-adishjain/golang-jwt/internal/config"
	helper "github.com/1107-adishjain/golang-jwt/internal/helpers"
	"github.com/1107-adishjain/golang-jwt/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GoogleLogin(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		url, codeVerifier, err := helper.GetGoogleOAuthURL(config.LoadConfig())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate PKCE code verifier"})
			return
		}
		// Store code_verifier in a secure cookie for the callback
		c.SetCookie("code_verifier", codeVerifier, 300, "", "/", true, true)
		c.Redirect(http.StatusFound, url)
	}
}

func GoogleCallback(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// after the user consents we will get the code from the query params
		code := c.Query("code")
		codeVerifier, err := c.Cookie("code_verifier")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing code_verifier for PKCE"})
			return
		}
		access_token, _, _, err := helper.ExchangeCodeForTokens(code, codeVerifier, config.LoadConfig())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to exchange code for tokens"})
			return
		}

		// then after this we will usee the access token to get the user info from google
		userInfo, err := helper.GetUserInfoFromGoogle(access_token)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user info from google"})
			return
		}

		// now get the user details from user Info map[string]interface{}
		email := userInfo.Email
		firstName := userInfo.GivenName
		lastName := userInfo.FamilyName

		// check if the user already exists in the database
		var user models.User
		if err := db.Where("email = ?", email).First(&user).Error; err != nil {
			// if user does not exist create a new user
			user = models.User{
				Email:     email,
				FirstName: firstName,
				LastName:  lastName,
				UserType:  "USER",
			}
			if err := db.Create(&user).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
				return
			}
		}

		// now generate a jwt token for the user
		accessToken, refreshToken, err := helper.GenerateJWT(user.UserID, user.UserType, user.FirstName, user.LastName)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate jwt token"})
			return
		}
		c.SetCookie(
			"refresh_token",
			refreshToken,
			7*24*60,
			"",
			"/",
			false,
			true,
		)

		c.JSON(http.StatusOK, gin.H{"access_token": accessToken, "user_id": user.UserID, "user_type": user.UserType, "first_name": user.FirstName, "last_name": user.LastName})
	}
}
