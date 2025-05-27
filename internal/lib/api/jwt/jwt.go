package jwt

import (
	"crypto/rand"
	"encoding/base64"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey []byte

func SetSecretKey(key string) {
	jwtKey = []byte(key)
}

func GetSecretKey() []byte {
	return jwtKey
}

func GenerateJWT(userID int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    "user",
		"exp":     time.Now().Add(15 * time.Minute).Unix(), // 15 * time.Minute
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func GenerateRefreshToken() (rawToken string, hashedToken string, err error) {
	bytes := make([]byte, 32)
	_, err = rand.Read(bytes)
	if err != nil {
		return "", "", err
	}
	rawToken = base64.URLEncoding.EncodeToString(bytes)

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(rawToken), bcrypt.DefaultCost)
	if err != nil {
		return "", "", err
	}
	hashedToken = string(hashedBytes)

	return rawToken, hashedToken, nil
}
