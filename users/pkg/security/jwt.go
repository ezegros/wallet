package security

import (
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type JWTData struct {
	UserID string
	jwt.StandardClaims
}

var secret = os.Getenv("JWT_KEY")

func CreateJWT(userID string) (string, error) {
	// Create JWT token
	claims := &JWTData{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: jwt.TimeFunc().Add(time.Hour * 24 * 7).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}

func ValidateJWT(tokenString string) (*JWTData, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTData{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTData); ok && token.Valid {
		return claims, nil
	}

	return nil, err
}
