package controller

import (
	"net/http"
	"strconv"
	"strings"

	model "github.com/Project-Sprint-Golang/EniQilo-Store/app/models"
	"github.com/Project-Sprint-Golang/EniQilo-Store/config"
	"github.com/Project-Sprint-Golang/EniQilo-Store/helper"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func RegisterStaff(c *gin.Context) {
	var user model.UserRegisterRequest
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Check if the phone number contains spaces
	if strings.Contains(user.PhoneNumber, " ") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Phone number cannot contain spaces"})
		return
	}
	// Check if the phone number starts with '+'
	if !strings.HasPrefix(user.PhoneNumber, "+") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Phone number must start with '+'"})
		return
	}

	hashedPassword, err := helper.GeneratePassword(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid request body"})
		return
	}

	// Check if the phone number already exists
	var count int
	err = config.DB.QueryRow("SELECT COUNT(*) FROM users WHERE phoneNumber = $1", user.PhoneNumber).Scan(&count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking phone number"})
		return
	}
	if count > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "Phone number already exists"})
		return
	}

	//setting role
	role := 2

	//insert user
	var lastInsertedID int
	err = config.DB.QueryRow("INSERT INTO users (phoneNumber, name, password, role) VALUES ($1, $2, $3, $4) RETURNING id",
		user.PhoneNumber, user.Name, string(hashedPassword), role).Scan(&lastInsertedID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error registering user"})
		return
	}

	token, err := helper.GenerateJWT(lastInsertedID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
		return
	}
	response := model.UserRegisterResponse{
		UserId:      strconv.Itoa(lastInsertedID),
		PhoneNumber: user.PhoneNumber,
		Name:        user.Name,
		AccessToken: token,
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"data":    response,
	})
}

func RegisterCustomer(c *gin.Context) {
	var user model.UserRegisterRequest
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Check if the phone number contains spaces
	if strings.Contains(user.PhoneNumber, " ") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Phone number cannot contain spaces"})
		return
	}
	// Check if the phone number starts with '+'
	if !strings.HasPrefix(user.PhoneNumber, "+") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Phone number must start with '+'"})
		return
	}

	hashedPassword, err := helper.GeneratePassword(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid request body"})
		return
	}

	// Check if the phone number already exists
	var count int
	err = config.DB.QueryRow("SELECT COUNT(*) FROM users WHERE phoneNumber = $1", user.PhoneNumber).Scan(&count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking phone number"})
		return
	}
	if count > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "Phone number already exists"})
		return
	}

	role := 1

	var lastInsertedID int
	err = config.DB.QueryRow("INSERT INTO users (phoneNumber, name, password, role) VALUES ($1, $2, $3, $4) RETURNING id",
		user.PhoneNumber, user.Name, string(hashedPassword), role).Scan(&lastInsertedID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error registering user"})
		return
	}

	token, err := helper.GenerateJWT(lastInsertedID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
		return
	}
	response := model.UserRegisterResponse{
		UserId:      strconv.Itoa(lastInsertedID),
		PhoneNumber: user.PhoneNumber,
		Name:        user.Name,
		AccessToken: token,
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"data":    response,
	})
}

func Login(c *gin.Context) {

	var user model.UserLoginRequest
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	// Check if the phone number contains spaces
	if strings.Contains(user.PhoneNumber, " ") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Phone number cannot contain spaces"})
		return
	}
	// Check if the phone number starts with '+'
	if !strings.HasPrefix(user.PhoneNumber, "+") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Phone number must start with '+'"})
		return
	}

	var userID int
	var hashedPassword, phoneNumber, name string

	err := config.DB.QueryRow("SELECT id, password, phoneNumber, name FROM users WHERE phoneNumber = $1", user.PhoneNumber).Scan(&userID, &hashedPassword, &phoneNumber, &name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(user.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect password"})
		return
	}
	token, err := helper.GenerateJWT(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
		return
	}
	response := model.UserRegisterResponse{
		UserId:      strconv.Itoa(userID),
		PhoneNumber: phoneNumber,
		Name:        name,
		AccessToken: token,
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "User registered successfully",
		"data":    response,
	})

}

func GetUsers(c *gin.Context) {
	var params model.GetCustomerParams
	if err := c.BindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query parameters"})
		return
	}
	// Construct the SQL query
	query := "SELECT userId, phoneNumber, name FROM customers WHERE 1=1 AND deletedAt IS NULL"
	args := []interface{}{}

	if params.Name != "" {
		query += " AND lower(name) LIKE $" + strconv.Itoa(len(args)+1)
		args = append(args, "%"+strings.ToLower(params.Name)+"%")
	}
}
