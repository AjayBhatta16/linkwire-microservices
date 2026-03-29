package myfunction

import (
	"log"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

func GetJWTUsername(tokenString string) (string, error) {
	secretKey := os.Getenv("JWT_SECRET_KEY")

	parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrTokenUnverifiable
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		log.Println("GetJWTUsername - Error parsing token:", err)
		return "", err
	}

	if !parsedToken.Valid {
		log.Println("GetJWTUsername - Invalid token")
		return "", jwt.ErrTokenInvalidClaims
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)

	if !ok {
		log.Println("GetJWTUsername - Error extracting claims")
		return "", jwt.ErrTokenInvalidClaims
	}

	username, ok := claims["username"].(string)

	if !ok || username == "" {
		log.Println("GetJWTUsername - Username claim missing or invalid")
		return "", jwt.ErrTokenInvalidClaims
	}

	return username, nil
}