package helper

import (
	"os"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

func GeneratePassword(password string) (string, error) {
	salt := os.Getenv("BCRYPT_SALT")
	saltInt, _ := strconv.Atoi(salt)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), int(saltInt))
	return string(hashedPassword), err
}
