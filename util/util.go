package util

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/denver-code/moza-backend/database"
	"github.com/denver-code/moza-backend/database/model"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// Helper functions for generating numbers and references
func GenerateAccountNumber() string {
	return fmt.Sprintf("%d", rand.Intn(900000000)+100000000)
}

func GenerateCardNumber() string {
	return fmt.Sprintf("%d", rand.Intn(9000000000000000)+1000000000000000)
}

func GenerateCVV() string {
	return fmt.Sprintf("%d", rand.Intn(900)+100)
}

func GenerateTransactionReference() string {
	return fmt.Sprintf("TXN%d", time.Now().UnixNano())
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPasswordHash compare password with hash
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
func ValidToken(t *jwt.Token, id string) bool {
	n, err := strconv.Atoi(id)
	if err != nil {
		return false
	}

	claims := t.Claims.(jwt.MapClaims)
	uid := int(claims["user_id"].(float64))

	return uid == n
}

func ValidUser(id string, p string) bool {
	db := database.DB
	var user model.User
	db.First(&user, id)
	if user.Username == "" {
		return false
	}
	if !CheckPasswordHash(p, user.Password) {
		return false
	}
	return true
}
