package middleware

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	Username string `json:"nim"`
	jwt.StandardClaims
}

func GenerateToken(username string, secretKey []byte) (string, error) {
	// Durasi token berlaku
	expirationTime := time.Now().Add(24 * time.Hour)

	// Membuat klaim JWT
	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	// Membuat token JWT dengan metode HMAC
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Menandatangani token dengan kunci rahasia
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyToken(tokenString string, secretKey []byte) (string, error) {
	// Parsing token dengan secret key
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return "", err
	}

	// Memeriksa apakah token valid
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims.Username, nil
	} else {
		return "", errors.New("Invalid token")
	}
}
