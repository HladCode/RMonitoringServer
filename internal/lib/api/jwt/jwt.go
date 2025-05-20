package jwt

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"time"

	"github.com/golang-jwt/jwt"
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
		"exp":     time.Now().Add(15 * time.Minute).Unix(),
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

	hash := sha256.Sum256([]byte(rawToken))
	hashedToken = base64.URLEncoding.EncodeToString(hash[:])

	return rawToken, hashedToken, nil
}
