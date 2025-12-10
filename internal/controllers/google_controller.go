package controllers

import (
	"net/http"
	"github.com/1107-adishjain/golang-jwt/internal/config"
	helper "github.com/1107-adishjain/golang-jwt/internal/helpers"
	"github.com/1107-adishjain/golang-jwt/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GoogleLogin(db *gorm.DB) gin.HandlerFunc{
	return func(c *gin.Context){
		// we will first redirect to the google oauth consent screen
		c.JSON(http.StatusOK, gin.H{"message": "Redirecting to Google consent screen"})

		url:= helper.GetGoogleOAuthURL(config.LoadConfig())

		// after we get the url we redirect the user to the google oauth consent screen
		c.Redirect(http.StatusFound, url)
	}  
}


func GoogleCallback(db *gorm.DB) gin.HandlerFunc{
	return func(c *gin.Context){
		// after the user consents we will get the code from the query params
		code:= c.Query("code")

		// now we should exchange the code for the access token and id token
		// now we will create a helper function which will exchange the code for the tokens
		access_token, _,_ ,err:= helper.ExchangeCodeForTokens(code, config.LoadConfig())

		if err!= nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to exchange code for tokens"})
			return
		}

		// then after this we will usee the access token to get the user info from google
		userInfo, err:= helper.GetUserInfoFromGoogle(access_token)
		if err!= nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user info from google"})
			return
		}

		// now get the user details from user Info map[string]interface{}
		email:= userInfo["email"].(string)
		firstName:= userInfo["given_name"].(string)
		lastName:= userInfo["family_name"].(string)
		// check if the user already exists in the database
		var user models.User		
		if err:= db.Where("email = ?", email).First(&user).Error; err!= nil{
			// if user does not exist create a new user
			user= models.User{
				Email: email,
				FirstName: firstName,
				LastName: lastName,
				UserType: "USER",
			}
			if err:= db.Create(&user).Error; err!= nil{
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
				return
			}
		}
	}
}



// entire code revision for google oauth login flow completed!!!!