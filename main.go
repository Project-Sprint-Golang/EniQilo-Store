package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Project-Sprint-Golang/EniQilo-Store/app/routes"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB

type User struct {
	UserID      int       `json:"id"`
	PhoneNumber string    `json:"phoneNumber"`
	Name        string    `json:"name"`
	Password    string    `json:"-"`
	Role        int       `json:"role"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type JWTClaims struct {
	UserID int `json:"userId"`
	jwt.StandardClaims
}

func main() {
	godotenv.Load()
	initDB()
	defer db.Close()

	router := gin.Default()

	// Setup routes
	routes.SetupRouter(router)

	router.Run(":8080")
}

func initDB() {
	connString := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?%s",
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PARAMS"))

	var err error
	db, err = sql.Open("postgres", connString)
	if err != nil {
		log.Fatal("Error connecting to the database: ", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Error pinging the database: ", err)
	}
}

func registerUser(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	user.Role = 2

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing the password"})
		return
	}

	_, err = db.Exec("INSERT INTO users (phoneNumber, name, password, role) VALUES ($1, $2, $3, $4)",
		user.PhoneNumber, user.Name, string(hashedPassword), user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error registering user"})
		return
	}

	var lastInsertedID int
	err = db.QueryRow("SELECT lastval()").Scan(&lastInsertedID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting last inserted ID"})
		return
	}

	token, err := generateJWT(lastInsertedID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"data": gin.H{
			"userID":      lastInsertedID,
			"phoneNumber": user.PhoneNumber,
			"name":        user.Name,
			"accessToken": token,
		},
	})
}

func registerStaff(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	user.Role = 1

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing the password"})
		return
	}

	_, err = db.Exec("INSERT INTO users (phoneNumber, name, password, role) VALUES ($1, $2, $3, $4)",
		user.PhoneNumber, user.Name, string(hashedPassword), user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error registering user"})
		return
	}

	var lastInsertedID int
	err = db.QueryRow("SELECT lastval()").Scan(&lastInsertedID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting last inserted ID"})
		return
	}

	token, err := generateJWT(lastInsertedID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Staff registered successfully",
		"data": gin.H{
			"userID":      lastInsertedID,
			"phoneNumber": user.PhoneNumber,
			"name":        user.Name,
			"accessToken": token,
		},
	})
}

func loginUser(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	var userID int
	var hashedPassword, phoneNumber, name string

	row := db.QueryRow("SELECT id, password, phoneNumber, name FROM users WHERE phoneNumber = $1", user.PhoneNumber)
	err := row.Scan(&userID, &hashedPassword, &phoneNumber, &name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(user.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect password"})
		return
	}

	token, err := generateJWT(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
		return
	}

	response := gin.H{
		"message": "User logged in successfully",
		"data": gin.H{
			"userId":      userID,
			"phoneNumber": phoneNumber,
			"name":        name,
			"accessToken": token,
		},
	}

	c.JSON(http.StatusOK, response)
}

func generateJWT(userID int) (string, error) {
	claims := JWTClaims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}
