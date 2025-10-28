package controllers

import (
	helper "github.com/1107-adishjain/golang-jwt/internal/helpers"
	"github.com/1107-adishjain/golang-jwt/internal/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
)

// this takes the password text/string and returns the hashed password using bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(bytes), err
}

func VerifyPassword(plainPassword, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
}

// in SignUp we take the email , password , usertype,first name ,lastname from the request body and then we check whether the email is already present in the database if not then we hash the password using the HashPassword function and then we create a new user in the database with the given details
func SignUp(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Email     string `json:"email"`
			Password  string `json:"password"`
			UserType  string `json:"user_type"`
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
		}
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}
		var existingUser models.User
		if err := db.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "email already in use"})
			return
		}
		// after it is confirmed that the email is not present in the database we hash the password using the HashPassword function
		hashedPassword, err := HashPassword(req.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
			return
		}
		user := models.User{
			Email:     req.Email,
			Password:  hashedPassword,
			UserType:  req.UserType,
			FirstName: req.FirstName,
			LastName:  req.LastName,
		}
		if err := db.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"message": "user created successfully"})
	}
}

// this handler will be used to login the user by taking the email and password from the request body and then we will verify the password using the VerifyPassword function and if the password is correct then we will generate a JWT token for the user and return it in the response
func Login(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}
		var user models.User
		if err := db.Where("email=?", req.Email).First(&user).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email"})
			return
		}
		if err := VerifyPassword(req.Password, user.Password); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid password"})
			return
		}
		// after the password is verified we will generate a jwt token for the user which will be used to authenticate the user in future requests
		access_token, refresh_token, err := helper.GenerateJWT(user.UserID, user.UserType, user.FirstName, user.LastName)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
			return
		}
		// this will set the refresh token in the http only cookie
		c.SetCookie(
			"refresh_token",
			refresh_token,
			7*24*60*60,
			"/",
			"",
			true,
			true,
		)
		c.JSON(http.StatusOK, gin.H{"access_token": access_token, "user_id": user.UserID, "user_type": user.UserType, "first_name": user.FirstName, "last_name": user.LastName})
	}
}

// get users will return the total number of users present in the database
func GetUsers(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var users []models.User
		res, err := db.Find(&users).RowsAffected, db.Find(&users).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		}
		if res == 0 {
			c.JSON(http.StatusOK, gin.H{"message": "no users found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"users": users})

	}
}

// this is use to get the specific user first the user id is validated then we use the user model struct to get the user from the database
func GetUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("id") //ID taken input from the user see the userRoutes.go file for it
		if err := helper.ValidateUserId(c, userId); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var user models.User
		if err := db.Where("user_id = ?", userId).First(&user).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"user": user})
	}
}
