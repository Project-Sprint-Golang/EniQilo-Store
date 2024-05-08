package main

import (
	"database/sql"
	"encoding/json"
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
	// Set up database connection
	godotenv.Load()
	initDB()
	defer db.Close()

	router := gin.Default()

	routes.SetupRouter(router)

	router.Run(":8080")

	// // Define HTTP routes
	// http.HandleFunc("/v1/staff/register", registerStaff)
	// http.HandleFunc("/v1/customer/register", registerUser)
	// http.HandleFunc("/v1/staff/login", loginUser)

	// // Start server
	// port := "8080"
	// log.Printf("Server listening on port %s...", port)
	// log.Fatal(http.ListenAndServe(":"+port, nil))
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

func registerUser(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Set role for regular user
	user.Role = 2

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error hashing the password", http.StatusInternalServerError)
		return
	}

	// Insert user into database
	_, err = db.Exec("INSERT INTO users (phoneNumber, name, password, role) VALUES ($1, $2, $3, $4)",
		user.PhoneNumber, user.Name, string(hashedPassword), user.Role)
	if err != nil {
		http.Error(w, "Error registering user", http.StatusInternalServerError)
		return
	}

	// Retrieve last inserted ID
	var lastInsertedID int
	err = db.QueryRow("SELECT lastval()").Scan(&lastInsertedID)
	if err != nil {
		http.Error(w, "Error getting last inserted ID", http.StatusInternalServerError)
		return
	}

	// Generate JWT token
	token, err := generateJWT(lastInsertedID)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "User registered successfully",
		"data": map[string]interface{}{
			"userID":      lastInsertedID,
			"phoneNumber": user.PhoneNumber,
			"name":        user.Name,
			"accessToken": token,
		},
	})
}

func registerStaff(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Set role for staff
	user.Role = 1

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error hashing the password", http.StatusInternalServerError)
		return
	}

	// Insert user into database
	_, err = db.Exec("INSERT INTO users (phoneNumber, name, password, role) VALUES ($1, $2, $3, $4)",
		user.PhoneNumber, user.Name, string(hashedPassword), user.Role)
	if err != nil {
		http.Error(w, "Error registering user", http.StatusInternalServerError)
		return
	}

	// Retrieve last inserted ID
	var lastInsertedID int
	err = db.QueryRow("SELECT lastval()").Scan(&lastInsertedID)
	if err != nil {
		http.Error(w, "Error getting last inserted ID", http.StatusInternalServerError)
		return
	}

	// Generate JWT token
	token, err := generateJWT(lastInsertedID)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Staff registered successfully",
		"data": map[string]interface{}{
			"userID":      lastInsertedID,
			"phoneNumber": user.PhoneNumber,
			"name":        user.Name,
			"accessToken": token,
		},
	})
}

func loginUser(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Declare variables to hold values from database
	var userID int
	var hashedPassword, phoneNumber, name string

	// Retrieve user from database
	row := db.QueryRow("SELECT id, password, phoneNumber, name FROM users WHERE phoneNumber = $1", user.PhoneNumber)
	err = row.Scan(&userID, &hashedPassword, &phoneNumber, &name)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Compare password
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(user.Password))
	if err != nil {
		http.Error(w, "Incorrect password", http.StatusBadRequest)
		return
	}

	// Generate JWT token
	token, err := generateJWT(userID)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	// Prepare the response
	response := map[string]interface{}{
		"message": "User logged in successfully",
		"data": map[string]interface{}{
			"userId":      userID,
			"phoneNumber": phoneNumber,
			"name":        name,
			"accessToken": token,
		},
	}

	// Return the response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
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
